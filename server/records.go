package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s Server) handlePostRecords(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("%s %s", "id", ErrValueShouldBeInt), http.StatusBadRequest)
		return
	}

	uRecord := &UserRecord{
		UserID: id,
		Record: Record{
			Name: "test",
			Type: "t1",
		},
	}

	if result := s.DB.Create(uRecord).Where(""); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	prepareResponseHeaders(w)

	log.Println(uRecord.UserID, uRecord.RecordID)
}

func (s Server) handleGetRecords(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("%s %s", "id", ErrValueShouldBeInt), http.StatusBadRequest)
		return
	}

	log.Println(id)

}
