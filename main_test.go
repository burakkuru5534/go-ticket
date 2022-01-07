// main_test.go
package main_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"log"

	_ "net/http"

)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tayitkan"
	dbname   = "ticketapp"
)


func TestMain(m *testing.M) {

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
	code := m.Run()
	os.Exit(code)
}

func ensureTableExists(db *sql.DB) {
	if _, err := db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
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

)`