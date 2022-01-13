// main_test.go
package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	_ "net/http"
)


func TestDbConnection(t *testing.T) {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	// close database
	defer db.Close()

	// check db
	err = db.Ping()

	fmt.Println("Connected!")

	ensureTableExists(db)

}

func ensureTableExists(db *sql.DB) {
	if _, err := db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func TestGetTicketOptions(t *testing.T) {

	req, err := http.NewRequest("GET", "/ticket_options", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	q := req.URL.Query()
	q.Add("id", "9368ec0b-701a-4be8-a267-b5deb1063128")
	req.URL.RawQuery = q.Encode()

	handler := http.HandlerFunc(GetTicketOption)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"ID":"9368ec0b-701a-4be8-a267-b5deb1063128","Name":"test ticket name","Desc":"test desc name","Allocation":1000,"CreatedAt":"2022-01-04T19:40:43.728787Z","UpdatedAt":"2022-01-04T19:40:43.728787Z"}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateTicketOptions(t *testing.T) {

	var jsonStr = []byte(`{
    "Name":"Test Ticket Options Name",
    "Desc":"There are 10.000 available tickets.",
    "Allocation":10000
}`)

	req, err := http.NewRequest("POST", "/ticket_options", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CreateTicketOption)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestPurchasesFromTicketOptions(t *testing.T) {

	wg := new(sync.WaitGroup)
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	// close database
	defer db.Close()

	// check db
	err = db.Ping()

	fmt.Println("Connected!")

	var jsonStr = []byte(`{
  "Quantity": 2,
  "UserID": "406c1d05-bbb2-4e94-b183-7d208c2692e1"
}`)

	req, err := http.NewRequest("POST", "/ticket_options", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "9368ec0b-701a-4be8-a267-b5deb1063128")
	req.URL.RawQuery = q.Encode()


	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := purchasesFromTicketOptions(wg,db)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}


const tableCreationQuery = `
--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.5
-- Dumped by pg_dump version 9.6.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: ar_internal_metadata; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE if not exists ar_internal_metadata (
                                                    key character varying NOT NULL,
                                                    value character varying,
                                                    created_at timestamp without time zone NOT NULL,
                                                    updated_at timestamp without time zone NOT NULL,

                                                    constraint ar_internal_metadata_pkey primary key (key)

);


--
-- Name: purchases; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE if not exists purchases (
                                         id uuid DEFAULT uuid_generate_v4() NOT NULL,
                                         quantity integer,
                                         user_id uuid,
                                         ticket_option_id uuid,
                                         created_at timestamp without time zone NOT NULL,
                                         updated_at timestamp without time zone NOT NULL,

                                         constraint purchases_pkey primary key (id)


);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE if not exists schema_migrations (
    version character varying NOT NULL,

    constraint schema_migrations_pkey primary key (version)

);


--
-- Name: ticket_options; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE if not exists ticket_options (
                                              id uuid DEFAULT uuid_generate_v4() NOT NULL,
                                              name character varying,
                                              "desc" character varying,
                                              allocation integer,
                                              created_at timestamp without time zone NOT NULL,
                                              updated_at timestamp without time zone NOT NULL,

                                              constraint ticket_options_pkey primary key (id)

);


--
-- Name: tickets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE if not exists tickets (
                                       id uuid DEFAULT uuid_generate_v4() NOT NULL,
                                       ticket_option_id uuid,
                                       purchase_id uuid,
                                       created_at timestamp without time zone NOT NULL,
                                       updated_at timestamp without time zone NOT NULL,
                                       constraint tickets_pkey primary key (id)

);

--
-- PostgreSQL database dump complete

`