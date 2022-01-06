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
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		r.Handle("/ticket/{id}/purchases", purchasesFromTicketOptions())
		wg.Done()
	}()
	go func() {
		r.Handle("/ticket/{id}", GetTicketOption())
		wg.Done()
	}()
	go func() {
		r.Handle("/ticket", CreateTicketOption())
		wg.Done()
	}()
	wg.Wait()

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

func purchasesFromTicketOptions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		// open database
		db, err := sql.Open("postgres", psqlconn)
		if err != nil {
			http.Error(w, "db connection error", 404)

		}
		// close database
		defer db.Close()

		// check db
		err = db.Ping()

		fmt.Println("Connected!")

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


	})
}
func GetTicketOption() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		var ticketOption models.TicketOptions

		params := mux.Vars(r)
		id := createKeyValuePairs(params)

		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		// open database
		db, err := sql.Open("postgres", psqlconn)
		if err != nil {
			http.Error(w, "db connection error", 404)
			return
		}
		// close database
		defer db.Close()

		// check db
		err = db.Ping()

		fmt.Println("Connected!")

		sq := fmt.Sprintf("select id::text,name,\"desc\",allocation,created_at,updated_at from  ticket_options where id::text = '%s'", id)
		err = db.QueryRow(sq).Scan(&ticketOption.ID, &ticketOption.Name, &ticketOption.Desc, &ticketOption.Allocation, &ticketOption.CreatedAt, &ticketOption.UpdatedAt)
		if err != nil {
			http.Error(w, "select from ticket options error", 404)
			return
		}

		json.NewEncoder(w).Encode(ticketOption)

	})
}
func CreateTicketOption() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		// open database
		db, err := sql.Open("postgres", psqlconn)
		if err != nil {
			http.Error(w, "db connection error", 404)
			return
		}
		// close database
		defer db.Close()

		// check db
		err = db.Ping()

		fmt.Println("Connected!")

		sq := fmt.Sprintf("insert into ticket_options (name, \"desc\", allocation, created_at, updated_at) values ('%s', '%s', %d, current_timestamp, current_timestamp)", ticketOption.Name.String, ticketOption.Desc.String, ticketOption.Allocation)
		_, err = db.Exec(sq)
		if err != nil {
			http.Error(w, "insert ticket options error", 404)
			return
		}

		json.NewEncoder(w).Encode(http.StatusOK)
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
