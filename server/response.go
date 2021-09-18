package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	ResponseStatusOK = iota
)

type Response struct {
	Status   int         `json:"status"`
	Response interface{} `json:"response,omitempty"`
}

func NewResponse(status int) *Response {
	return &Response{Status: status}
}

func (r *Response) WithResponse(i interface{}) *Response {
	r.Response = i
	return r
}

func (r *Response) Json() ([]byte, error) {
	j, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed while marshaling json: %v", err))
	}

	return j, err
}

func (r *Response) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	j, _ := r.Json()
	w.Write(j)
}
