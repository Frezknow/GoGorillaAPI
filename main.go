package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}
type allEvents []event

var events = allEvents{}
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:reynolds6721@tcp(localhost:3306)/gorilla_api")
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}
	db.SetMaxIdleConns(100)

	err = db.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	router := mux.NewRouter()
	//.StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/event/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/event/{id}", updateEvent).Methods("PUT")
	router.HandleFunc("/event/{id}", updateEvent).Methods("DELETE")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, err := db.Query("select id,title,description from events")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var newEvent event
		err := result.Scan(&newEvent.ID, &newEvent.Title, &newEvent.Description)
		if err != nil {
			panic(err.Error())
		}
		events = append(events, newEvent)
	}
	json.NewEncoder(w).Encode(events)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	//var newEvent event
	new, err := db.Prepare("INSERT INTO events(title,description) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	keyVal := make(map[string]string)
	json.Unmarshal(reqBody, &keyVal) //newEVent
	title := keyVal["title"]
	description := keyVal["description"]
	_, err = new.Exec(title, description)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Event was created successfully.")
}
func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update.")
	}
	json.Unmarshal(reqBody, &updatedEvent)
	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)

		}
	}
}
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}
func getOneEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	eventID := mux.Vars(r)["id"]
	result, err := db.Query("SELECT * FROM events WHERE id = ?", eventID)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var newEvent event
	for result.Next() {
		err := result.Scan(&newEvent.ID, &newEvent.Title, &newEvent.Description)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(newEvent)
	// for _, singleEvent := range events {
	// 	if singleEvent.ID == eventID {
	// 		json.NewEncoder(w).Encode(singleEvent)
	// 	}
	// }
}
func homeLink(q http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(q, "Welcome home!")
}

//docker related commands docker run --name goMysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=reynolds6721 mysql
