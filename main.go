package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var connStr = "user=postgres password=postgres dbname=timetracker sslmode=disable"

var db, fErr = sql.Open("postgres", connStr)

func main() {
	if fErr != nil {
		panic(fErr)
	}
	defer db.Close()
	r := mux.NewRouter()
	
	r.HandleFunc("/groups", groupsGetHandler).Methods("GET")
	r.HandleFunc("/tasks", tasksGetHandler).Methods("GET")

	r.HandleFunc("/groups", groupsPostHandler).Methods("POST")
	r.HandleFunc("/tasks", tasksPostHandler).Methods("POST")
	r.HandleFunc("/timeframes", timeframesPostHandler).Methods("POST")

	r.HandleFunc("/groups/{id:[0-9]+}", groupsPutHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id:[0-9]+}", tasksPutHandler).Methods("PUT")

	r.HandleFunc("/groups/{id:[0-9]+}", groupsDeleteHandler).Methods("DELETE")
	r.HandleFunc("/tasks/{id:[0-9]+}", tasksDeleteHandler).Methods("DELETE")
	r.HandleFunc("/timeframes/{id:[0-9]+}", timeframesDeleteHandler).Methods("DELETE")

	http.Handle("/", r)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Panic(err)
	}
}
