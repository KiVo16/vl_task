package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	ErrValueNotFound = iota
	ErrValueInvalidType
	ErrBodyMissing
	ErrBodyRead
	ErrJsonInvalid
	ErrInsertData
	ErrGetData
)

var errMap = map[int]string{
	ErrValueNotFound:    "Value not found",
	ErrValueInvalidType: "Expected %v instead got %v",
	ErrBodyMissing:      "Body is missing",
	ErrBodyRead:         "Body read error",
	ErrJsonInvalid:      "Invalid json",
	ErrInsertData:       "Error while inserting data",
	ErrGetData:          "Error while getting data",
}

type ServerError struct {
	HTTPCode      int    `json:"http_code"`
	ErrorCode     int    `json:"error_code"`
	RefersTo      string `json:"refers_to,omitempty"`
	Message       string `json:"message,omitempty"`
	DetailedError error  `json:"detailed_error,omitempty"`
}

func NewServerError(httpError, errorCode int) *ServerError {
	return &ServerError{ErrorCode: errorCode, HTTPCode: httpError}
}

func NewPredefinedServerError(httpError, errorCode int, args ...interface{}) *ServerError {
	log.Println(errMap[errorCode])
	return &ServerError{ErrorCode: errorCode, HTTPCode: httpError, Message: fmt.Sprintf(errMap[errorCode], args...)}
}

func (e *ServerError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.HTTPCode)
	j, _ := e.Json()
	w.Write(j)
}

func (e *ServerError) WithMessage(msg string) *ServerError {
	e.Message = msg
	return e
}

func (e *ServerError) WithRefersTo(refers string) *ServerError {
	e.RefersTo = refers
	return e
}

func (e *ServerError) WithDetailedError(err error) *ServerError {
	e.DetailedError = err
	return e
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("ErrorCode: %d, RefersTo: %s, Message: %s, DetailedError: %v", e.ErrorCode, e.RefersTo, e.Message, e.DetailedError)
}

func (e *ServerError) Json() ([]byte, error) {
	j, err := json.MarshalIndent(e, "", " ")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed while marshaling json: %v", err))
	}

	return j, err
}
