package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gorilla_api/models"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type user struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}
type event struct {
	ID          string    `json:"ID"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Comments    []comment `json:"comments"`
}
type comment struct {
	ID    string `json:"ID"`
	Body  string `json:"Body"`
	Owner int    `json:"Owner"`
}
type allEvents []event

var events = allEvents{}
var db *sql.DB
var db2 *gorm.DB

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
	router.HandleFunc("/comment", createComment).Methods("POST")
	router.HandleFunc("/comment/{id}", deleteComment).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// This is the root directory of uploaded files
var base = "/home/mehrdadep/example"

func Upload(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	n := fmt.Sprintf("%d - %s", time.Now().UTC().Unix(), file.Filename)
	dst := fmt.Sprintf("%s/%s", base, n)
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return n, err
}
func createComment(w http.ResponseWriter, r *http.Request) {
	//var newEvent event
	new, err := db.Prepare("INSERT INTO comments(owner,body,event_id) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter a comment for this event.")
	}
	keyVal := make(map[string]string)
	json.Unmarshal(reqBody, &keyVal)
	body := keyVal["body"]
	owner, err := strconv.ParseInt(keyVal["owner"], 0, 64)
	if err != nil {
		panic(err)
	}
	eventId, err := strconv.ParseInt(keyVal["event_id"], 0, 64)
	if err != nil {
		panic(err)
	}
	_, err = new.Exec(owner, body, eventId)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Comment was created successfully.")
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

func deleteComment(w http.ResponseWriter, r *http.Request) {
	// commentID := mux.Vars(r)["id"]
	// for i, singleComment := range events {
	// 	if singleComment.ID == commentID {
	// 		events = append(comments[:i], comments[i+1:]...)
	// 		fmt.Fprintf(w, "The comment with ID %v has been deleted successfully", commentID)
	// 	}
	// }
}
func getOneEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	eventID := mux.Vars(r)["id"]
	result, err := db.Query("SELECT * FROM events  WHERE events.id = ?", eventID)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	//var newEvent event
	// for result.Next() {
	// 	//[].Id, &newEvent.&Comments[].Body, &newEvent.Comments[].Owner
	// 	err := result.Scan(&newEvent.ID, &newEvent.Title, &newEvent.Description))
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// }
	//json.NewEncoder(w).Encode(newEvent)
	var eventModel models.EventModel
	events, _ := eventModel.FindByID(eventID)
	for _, event := range events {
		fmt.Println(event.ToString(), " \nSearched for ID:", eventID)
		fmt.Println("Comments: ", len(event.Comments))
		if len(event.Comments) > 0 {
			for _, comment := range event.Comments {
				fmt.Println(comment.ToString())
				fmt.Println("=============================")
			}
		}
		fmt.Println("--------------------")
	}
}

func homeLink(q http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(q, "Welcome home!")
}

//docker related commands docker run --name goMysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=reynolds6721 mysql
