package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s Server) handlePostRecords(w http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return
	}

	name, err := extractStringFromMap("name", m)
	if err != nil {
		return
	}

	t, err := extractStringFromMap("type", m)
	if err != nil {
		return
	}

	record, err := s.createRecord(name, t)
	if err != nil {
		return
	}
	prepareResponseHeaders(w)

	log.Println(record.ID)
}

func (s Server) handleAssignRecordToUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("%s %s", "id", ErrValueShouldBeInt), http.StatusBadRequest)
		return
	}

	recordID, err := strconv.Atoi(params["recordID"])
	if err != nil {
		http.Error(w, fmt.Sprintf("%s %s", "recordID", ErrValueShouldBeInt), http.StatusBadRequest)
		return
	}

	if err := s.assignRecordToUser(id, recordID); err != nil {
		return
	}

}

func (s Server) handleCountRecords(w http.ResponseWriter, req *http.Request) {
	var count int64

	q := s.db.Table("user_records").
		Joins("INNER JOIN records AS r ON r.id = record_id").
		Joins("INNER JOIN users AS u ON u.id = user_id")

	if len(req.FormValue("type")) > 0 {
		q.Where("r.type = ?", req.FormValue("type"))
	}

	if len(req.FormValue("user_name")) > 0 {
		q.Where("u.name = ?", req.FormValue("user_name"))
	}

	if result := q.Count(&count); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(count)
}

func (s Server) createRecord(name, t string) (*Record, error) {
	record := &Record{Name: name, Type: t}

	if result := s.db.Create(record); result.Error != nil {
		return nil, result.Error
	}

	return record, nil
}

func (s Server) assignRecordToUser(userID, recordID int) error {
	a := &UserRecord{
		UserID:   userID,
		RecordID: recordID,
	}

	if result := s.db.Create(a); result.Error != nil {
		return result.Error
	}

	return nil
}
