package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (s Server) handlePostUsers(w http.ResponseWriter, req *http.Request) {

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

	user, err := s.createUser(name)
	if err != nil {
		return
	}

	log.Println(user.ID)

	prepareResponseHeaders(w)
}

func (s Server) handleGetUsers(w http.ResponseWriter, req *http.Request) {
	limit, offset := defaultLimit, defaultOffset

	if len(req.FormValue("limit")) > 0 {
		v, err := strconv.Atoi(req.FormValue("limit"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		limit = v
	}

	if len(req.FormValue("offset")) > 0 {
		v, err := strconv.Atoi(req.FormValue("offset"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		offset = v
	}

	users := []User{}
	if result := s.db.Limit(limit).Offset(offset).Find(&users); result.Error != nil {

		return
	}

	j, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		return
	}

	w.Write(j)
}

func (s Server) createUser(name string) (*User, error) {
	user := &User{Name: name}

	if result := s.db.Create(user); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
