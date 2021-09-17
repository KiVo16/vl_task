package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	version       = 1.0
	defaultLimit  = 10
	defaultOffset = -1

	ErrValueShouldBeInt = "must be of type integer"
)

var defaultResponseHeaders = map[string]string{
	"Content-Type": "application/json",
}

type Server struct {
	DB *gorm.DB
}

func main() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed while connecting to database: %v", err)
	}

	db.AutoMigrate(&User{}, &Record{}, &UserRecord{})

	server := Server{DB: db}
	m := mux.NewRouter()

	u := m.PathPrefix("/{users:users\\/?}").Subrouter()
	u.HandleFunc("/", server.handlePostUsers).Methods("POST")
	u.HandleFunc("/", server.handleGetUsers).Methods("GET").Queries("limit", "{limit}", "offset", "{offset}")

	r := u.PathPrefix("/{/{id}/records:/{id}/records\\/?}").Subrouter()
	r.HandleFunc("/", server.handleGetRecords).Methods("GET")
	r.HandleFunc("/", server.handlePostRecords).Methods("POST")

	if err := http.ListenAndServe(":8000", m); err != nil {
		log.Fatalf("Failed while starting server: %v", err)
	}

	log.Println("Server started")
}
