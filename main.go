package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}
type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for chance to learn how golang works and get to eventually try it out",
	},
	{
		ID:          "2",
		Title:       "Introduction to Golang pt2",
		Description: "Come join us for chance to learn how golang works and get to eventually try it out",
	},
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
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
func deleteEvent(w http.ResponseWriter, r *http.Request){
	eventID := mux.Vars(r)["id"]
	for i, singleEvent := range events{
		if singleEvent.ID == eventID{
		events = append(events[:i], events[i+1:]...)
		fmt.Fprintf(w, "The event with ID %v has been deleted successfully",eventID)
	    }
	}
}
func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}
func homeLink(q http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(q, "Welcome home!")
}
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/event/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/event/{id}", updateEvent).Methods("PUT")
	router.HandleFunc("/event/{id}", updateEvent).Methods("DELETE")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
