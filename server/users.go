package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"gorm.io/gorm"
)

func (s Server) handlePostUsers(w http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		NewPredefinedServerError(http.StatusBadRequest, ErrBodyRead).WithRefersTo("body").WithDetailedError(err).Write(w)
		return
	}

	if len(body) == 0 {
		NewPredefinedServerError(http.StatusBadRequest, ErrBodyMissing).WithRefersTo("body").Write(w)
		return
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		NewPredefinedServerError(http.StatusBadRequest, ErrJsonInvalid).WithRefersTo("body").WithDetailedError(err).Write(w)
		return
	}

	name, err := extractStringFromMap("name", m)
	if err != nil {
		handleExtractStringFromMapError(w, "name", err)
		return
	}

	user, err := s.createUser(name)
	if err != nil {
		NewPredefinedServerError(http.StatusInternalServerError, ErrInsertData).WithDetailedError(err).Write(w)
		return
	}

	NewResponse(ResponseStatusOK).WithResponse(struct {
		CreatedID int `json:"created_id"`
	}{user.ID}).Write(w)
}

func (s Server) handleGetUsers(w http.ResponseWriter, req *http.Request) {
	limit, offset := defaultLimit, defaultOffset

	if len(req.FormValue("limit")) > 0 {
		v, err := strconv.Atoi(req.FormValue("limit"))
		if err != nil {
			NewPredefinedServerError(http.StatusBadRequest, ErrValueInvalidType, "int", reflect.TypeOf(v)).WithRefersTo("limit").WithDetailedError(err).Write(w)
			return
		}

		limit = v
	}

	if len(req.FormValue("offset")) > 0 {
		v, err := strconv.Atoi(req.FormValue("offset"))
		if err != nil {
			NewPredefinedServerError(http.StatusBadRequest, ErrValueInvalidType, "int", reflect.TypeOf(v)).WithRefersTo("offset").WithDetailedError(err).Write(w)
			return
		}

		offset = v
	}

	users := []User{}
	if err := s.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return
		}
		NewPredefinedServerError(http.StatusInternalServerError, ErrGetData).WithDetailedError(err).Write(w)
		return
	}

	NewResponse(ResponseStatusOK).WithResponse(users).Write(w)
}

func (s Server) createUser(name string) (*User, error) {
	user := &User{Name: name}

	if result := s.db.Create(user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
