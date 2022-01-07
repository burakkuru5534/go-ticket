package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"example.com/m/models"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	uuid2 "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"time"

	"github.com/rs/cors"

	_ "database/sql"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tayitkan"
	dbname   = "ticketapp"
)

func main() {

	//router
	r := mux.NewRouter()
	//api endpoints

	defer timeTrack(time.Now(), "purchase process info")
	fmt.Println("Starting concurrent calls...")

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


	initDB(db)
	wg := new(sync.WaitGroup)

		//r.Handle("/ticket/{id}/purchases", purchasesFromTicketOptions(wg))
	go r.Handle("/ticket_options/{id}/purchases", purchasesFromTicketOptions(wg, db))

		r.Handle("/ticket_options/{id}", GetTicketOption(wg, db))

		r.Handle("/ticket_options", CreateTicketOption(wg, db))


	//define options
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})
	//start server
	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(r)))

}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func purchasesFromTicketOptions(wg *sync.WaitGroup, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		wg.Wait()
		wg.Add(1)

		w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		ticketOptionID := mux.Vars(r)
		ticketOptionsID := createKeyValuePairs(ticketOptionID)

		var purchases models.Purchases
		err := BodyToJsonReq(r, &purchases)
		if err != nil {
			http.Error(w, "body to json request error", 404)

		}

		ticketopt, err := uuid.FromString(ticketOptionsID)
		purchases.TicketOptionID = uuid2.UUID(ticketopt)



		allocation := getAllocationOfTicketOptions(uuid.UUID(purchases.TicketOptionID), db)

		isThereAvailableTickets := quantityAndAllocationCompare(purchases.Quantity.Int64, allocation)

		if isThereAvailableTickets {

			var purchaseID string
			sq := fmt.Sprintf("insert into purchases (quantity, user_id, ticket_option_id, created_at, updated_at) values (%d, '%v', '%v', current_timestamp, current_timestamp) returning id", purchases.Quantity.Int64, purchases.UserID, ticketOptionsID)
			err = db.QueryRow(sq).Scan(&purchaseID)
			if err != nil {
				http.Error(w, "insert purchases error", 404)

			}

			sq = fmt.Sprintf("insert into tickets (ticket_option_id, purchase_id, created_at, updated_at) values ('%v', '%v', current_timestamp, current_timestamp) ", ticketOptionsID, purchaseID)
			_, err = db.Exec(sq)
			if err != nil {
				http.Error(w, "insert tickets error", 404)

			}

			isAllocationDescreased := decreaseAllocationOfTicket(ticketOptionsID, purchases.Quantity.Int64, db)
			if !isAllocationDescreased {
				http.Error(w, " allocation decrease error", 404)

			}
			json.NewEncoder(w).Encode(http.StatusOK)

		}else {
			http.Error(w, "there is not any available tickets", 404)
		}

		wg.Done()

	})
}
func GetTicketOption(wg *sync.WaitGroup, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Wait()
		wg.Add(1)
		w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		var ticketOption models.TicketOptions

		params := mux.Vars(r)
		id := createKeyValuePairs(params)


		sq := fmt.Sprintf("select id::text,name,\"desc\",allocation,created_at,updated_at from  ticket_options where id::text = '%s'", id)
		err := db.QueryRow(sq).Scan(&ticketOption.ID, &ticketOption.Name, &ticketOption.Desc, &ticketOption.Allocation, &ticketOption.CreatedAt, &ticketOption.UpdatedAt)
		if err != nil {
			http.Error(w, "select from ticket options error", 404)
			return
		}

		json.NewEncoder(w).Encode(ticketOption)
		wg.Done()

	})

}
func CreateTicketOption(wg *sync.WaitGroup, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Wait()
		wg.Add(1)
		w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		var ticketOption models.TicketOptions
		err := BodyToJsonReq(r, &ticketOption)
		if err != nil {
			http.Error(w, "body to json request error", 404)
			return
		}


		sq := fmt.Sprintf("insert into ticket_options (name, \"desc\", allocation, created_at, updated_at) values ('%s', '%s', %d, current_timestamp, current_timestamp)", ticketOption.Name.String, ticketOption.Desc.String, ticketOption.Allocation)
		_, err = db.Exec(sq)
		if err != nil {
			http.Error(w, "insert ticket options error", 404)
			return
		}

		json.NewEncoder(w).Encode(http.StatusOK)
		wg.Done()
	})
}

//helper functions
func BodyToJsonReq(r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return nil
}
func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for _, value := range m {
		fmt.Fprintf(b, "%s", value)
	}
	return b.String()
}
func quantityAndAllocationCompare(quantity int64, allocation int64) bool {

	return allocation >= quantity
}
func decreaseAllocationOfTicket(ticketOptionsID string, quantity int64, db *sql.DB) bool {

	sq := fmt.Sprintf("update ticket_options set allocation = allocation - %d where id::text = '%s'", quantity, ticketOptionsID)
	_, err := db.Exec(sq)
	if err != nil {
		return false
	}

	return true
}
func getAllocationOfTicketOptions(ticketOptionsID uuid.UUID, db *sql.DB) int64 {
	var allocation int64
	sq := fmt.Sprintf("select coalesce(allocation,0) from ticket_options where id::text = '%s'", ticketOptionsID)
	err := db.QueryRow(sq).Scan(&allocation)
	if err != nil {
		return 0
	}

	return allocation
}

func initDB (db *sql.DB){

	sq := `
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
--
`

	_, err := db.Exec(sq)
	if err != nil {

	}

}
