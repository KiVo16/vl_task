package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s Server) handlePostUsers(w http.ResponseWriter, req *http.Request) {

	user := &User{Name: "Jakub"}
	if result := s.DB.Create(user); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(user.ID)

	prepareResponseHeaders(w)
}

func (s Server) handleGetUsers(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	log.Println(params)
	limit, offset := defaultLimit, defaultOffset

	if val, ok := params["limit"]; ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		limit = v
	}

	if val, ok := params["offset"]; ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		offset = v
	}

	users := []User{}
	if result := s.DB.Limit(limit).Offset(offset).Find(&users); result.Error != nil {

		return
	}

	j, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		return
	}

	w.Write(j)
}
