package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	sampleDataFlag := flag.Bool("load-sample-data", false, "Loads sample data. Default value: false")
	dbPath := flag.String("db", "data.db", "Path to SQLite database file. New database will be created if provided path points to non-existing file. Default value: ./data.db")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed while connecting to database: %v", err)
	}

	db.Exec("PRAGMA foreign_keys = ON")
	db.AutoMigrate(&User{}, &Record{}, &UserRecord{})

	server := Server{db: db}
	if *sampleDataFlag == true {
		server.loadTestData()
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
