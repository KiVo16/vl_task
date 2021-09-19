package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
)

var defaultResponseHeaders = map[string]string{
	"Content-Type": "application/json",
}

type Server struct {
	m   *mux.Router
	srv *http.Server
	db  *gorm.DB
}

func (s *Server) Init(dbPath string) {
	rand.Seed(time.Now().UnixNano())

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed while connecting to database: %v", err)
	}

	db.Exec("PRAGMA foreign_keys = ON")
	db.AutoMigrate(&User{}, &Record{}, &UserRecord{})

	s.db = db
	s.m = mux.NewRouter()

}

func (s *Server) Run(ssl bool) {
	s.srv = &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.m,
	}

	s.m.HandleFunc("/{users:users\\/?}", s.handlePostUsers).Methods("POST")
	s.m.HandleFunc("/{users:users\\/?}", s.handleGetUsers).Methods("GET")

	s.m.HandleFunc("/{records:records\\/?}", s.handlePostRecords).Methods("POST")
	s.m.HandleFunc("/{records:records\\/?}", s.handleCountRecords).Methods("GET")

	s.m.HandleFunc("/users/{id}/records/{recordID}", s.handleAssignRecordToUser).Methods("POST")

	go func() {
		var err error
		// err = s.srv.ListenAndServeTLS("../cert.crt", "../cert.key")

		err = s.srv.ListenAndServe()

		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Server started")

}

func main() {
	sampleDataFlag := flag.Bool("load-sample-data", false, "Loads sample data.")
	sampleDataNamesPathFlag := flag.String("sample-users-names", "./sampleData/sampleNames.json", "Path to sample data for user names creation. JSON must contains only array of single string values.")
	sampleDataRecordsPathFlag := flag.String("sample-records-names", "./sampleData/sampleRecords.json", "Path to sample data for record names creation. JSON must contains only array of single string values.")
	dbPath := flag.String("db", "data.db", "Path to SQLite database file. New database will be created if provided path points to non-existing file.")
	flag.Parse()

	server := Server{}
	server.Init(*dbPath)

	if *sampleDataFlag == true {
		if err := server.loadTestData(*sampleDataNamesPathFlag, *sampleDataRecordsPathFlag); err != nil {
			log.Fatalf("Failed while loading sample data: %v", err)
		}
	}

	server.Run(true)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server gracefull shutdown failed: %v", err)
	}

	log.Println("Server is shutting down")
	os.Exit(0)
}
