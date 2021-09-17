package main

import (
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

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
	db *gorm.DB
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed while connecting to database: %v", err)
	}

	db.Exec("PRAGMA foreign_keys = ON")
	db.AutoMigrate(&User{}, &Record{}, &UserRecord{})

	server := Server{db: db}
	err = server.loadTestData()
	if err != nil {
		log.Fatalf("Failed while loading sample data: %v", err)
	}
	m := mux.NewRouter()

	m.HandleFunc("/{users:users\\/?}", server.handlePostUsers).Methods("POST")
	m.HandleFunc("/{users:users\\/?}", server.handleGetUsers).Methods("GET")

	m.HandleFunc("/{records:records\\/?}", server.handlePostRecords).Methods("POST")
	m.HandleFunc("/{records:records\\/?}", server.handleCountRecords).Methods("GET")

	m.HandleFunc("/users/{id}/records/{recordID}", server.handleAssignRecordToUser).Methods("POST")

	if err := http.ListenAndServe(":8000", m); err != nil {
		log.Fatalf("Failed while starting server: %v", err)
	}

	log.Println("Server started")

}
