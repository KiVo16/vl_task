package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s Server) handlePostRecords(w http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		NewPredefinedServerError(http.StatusBadRequest, ErrBodyMissing).Write(w)
		return
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return
	}

	name, err := extractStringFromMap("name", m)
	if err != nil {
		handleExtractStringFromMapError(w, "name", err)
		return
	}

	t, err := extractStringFromMap("type", m)
	if err != nil {
		handleExtractStringFromMapError(w, "type", err)
		return
	}

	record, err := s.createRecord(name, t)
	if err != nil {
		NewPredefinedServerError(http.StatusInternalServerError, ErrInsertData).WithDetailedError(err).Write(w)
		return
	}

	NewResponse(ResponseStatusOK).WithResponse(struct {
		CreatedID int `json:"created_id"`
	}{record.ID}).Write(w)
}

func (s Server) handleAssignRecordToUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		NewPredefinedServerError(http.StatusBadRequest, ErrValueInvalidType, "int").WithRefersTo("user_id").WithDetailedError(err).Write(w)
		return
	}

	recordID, err := strconv.Atoi(params["recordID"])
	if err != nil {
		NewPredefinedServerError(http.StatusBadRequest, ErrValueInvalidType, "int").WithRefersTo("record_id").WithDetailedError(err).Write(w)
		return
	}

	if err := s.assignRecordToUser(id, recordID); err != nil {
		if err.Error() == "FOREIGN KEY constraint failed" {
			// Nie wiem, którego kodu http użyć. Jest to błąd po stronie serwera spowodowany złymi danymi, które są przyczyną błędu
			// po wykonaniu requestu do bazy. Nie są błęde w samym zapytaniu więc http.StatusBadRequest raczej odpada.
			NewPredefinedServerError(http.StatusInternalServerError, ErrForeignKey).WithDetailedError(err).Write(w)
			return
		}

		NewPredefinedServerError(http.StatusInternalServerError, ErrInsertData).WithDetailedError(err).Write(w)
		return
	}

	NewResponse(ResponseStatusOK).WithResponse(struct {
		UserID   int `json:"user_id"`
		RecordID int `json:"record_id"`
	}{id, recordID}).Write(w)

}

func (s Server) handleCountRecords(w http.ResponseWriter, req *http.Request) {
	var count int64

	q := s.db.Table("user_records").
		Joins("INNER JOIN records AS r ON r.id = record_id").
		Joins("INNER JOIN users AS u ON u.id = user_id")

	if len(req.FormValue("type")) > 0 {
		q.Where("r.type = ?", req.FormValue("type"))
	} else {
		NewPredefinedServerError(http.StatusBadRequest, ErrValueNotFound).WithRefersTo("type").Write(w)
		return
	}

	if len(req.FormValue("user_name")) > 0 {
		q.Where("u.name = ?", req.FormValue("user_name"))
	} else {
		NewPredefinedServerError(http.StatusBadRequest, ErrValueNotFound).WithRefersTo("user_name").Write(w)
		return
	}

	if err := q.Count(&count).Error; err != nil {
		NewPredefinedServerError(http.StatusInternalServerError, ErrGetData).WithDetailedError(err).Write(w)
		return
	}

	NewResponse(ResponseStatusOK).WithResponse(struct {
		Count int64 `json:"count"`
	}{count}).Write(w)
}

func (s Server) createRecord(name, t string) (*Record, error) {
	record := &Record{Name: name, Type: t}

	if err := s.db.Create(record).Error; err != nil {
		return nil, err
	}

	return record, nil
}

func (s Server) assignRecordToUser(userID, recordID int) error {
	a := &UserRecord{
		UserID:   userID,
		RecordID: recordID,
	}

	if err := s.db.Create(a).Error; err != nil {
		return err
	}

	return nil
}
