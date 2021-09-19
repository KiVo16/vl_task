package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHandlePostRecords(t *testing.T) {

	tests := []RestTest{
		NewRestTest("/records", "POST", TestMap{"name": "test", "type": "t1"}, TestMap{"status": 0, "response": TestMap{"created_id": 4}}, http.StatusOK, false),
		NewRestTest("/records", "POST", TestMap{"name": "test", "type": "t1"}, TestMap{"status": 0, "response": TestMap{"created_id": 5}}, http.StatusOK, false),
		NewRestTest("/records", "POST", TestMap{"name": "test"}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueNotFound,
			"message":    errMap[ErrValueNotFound],
			"refers_to":  "type",
		}, http.StatusBadRequest, false),
		NewRestTest("/records", "POST", TestMap{"type": "t1"}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueNotFound,
			"message":    errMap[ErrValueNotFound],
			"refers_to":  "name",
		}, http.StatusBadRequest, false),
		NewRestTest("/records", "POST", TestMap{"name": 25, "type": "t1"}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    "Expected string",
			"refers_to":  "name",
		}, http.StatusBadRequest, false),
		NewRestTest("/records", "POST", TestMap{"name": "test", "type": true}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    "Expected string",
			"refers_to":  "type",
		}, http.StatusBadRequest, false),
	}

	for _, test := range tests {
		test.Test(t)
	}

}

func TestHandleGetRecords(t *testing.T) {

	tests := []RestTest{
		NewRestTest("/records?user_name=Dominika&type=t1", "GET", TestMap{}, TestMap{"status": 0, "response": TestMap{"count": 2}}, http.StatusOK, false),
		NewRestTest("/records?user_name=Dominika&type=t2", "GET", TestMap{}, TestMap{"status": 0, "response": TestMap{"count": 1}}, http.StatusOK, false),
		NewRestTest("/records?user_name=Dominika", "GET", TestMap{}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueNotFound,
			"message":    errMap[ErrValueNotFound],
			"refers_to":  "type",
		}, http.StatusBadRequest, false),
		NewRestTest("/records?type=t1", "GET", TestMap{}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueNotFound,
			"message":    errMap[ErrValueNotFound],
			"refers_to":  "user_name",
		}, http.StatusBadRequest, false),
	}

	for _, test := range tests {
		test.Test(t)
	}

}

func TestHandleAssignRecordToUser(t *testing.T) {

	tests := []RestTest{
		NewRestTest("/users/1/records/1", "POST", TestMap{}, TestMap{"status": 0, "response": TestMap{"user_id": 1, "record_id": 1}}, http.StatusOK, false),
		NewRestTest("/users/1/records/2", "POST", TestMap{}, TestMap{"status": 0, "response": TestMap{"user_id": 1, "record_id": 2}}, http.StatusOK, false),
		NewRestTest("/users/1/records/1", "POST", TestMap{}, TestMap{
			"http_code":  http.StatusInternalServerError,
			"error_code": ErrInsertData,
			"message":    errMap[ErrInsertData],
		}, http.StatusInternalServerError, false),
		NewRestTest("/users/s/records/1", "POST", TestMap{}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    fmt.Sprintf(errMap[ErrValueInvalidType], "int"),
			"refers_to":  "user_id",
		}, http.StatusBadRequest, false),
		NewRestTest("/users/1/records/s", "POST", TestMap{}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    fmt.Sprintf(errMap[ErrValueInvalidType], "int"),
			"refers_to":  "record_id",
		}, http.StatusBadRequest, false),
		NewRestTest("/users/1/records/", "POST", TestMap{}, TestMap{}, http.StatusNotFound, true),
		NewRestTest("/users/1/records/20", "POST", TestMap{}, TestMap{
			"http_code":  http.StatusInternalServerError,
			"error_code": ErrForeignKey,
			"message":    errMap[ErrForeignKey],
		}, http.StatusInternalServerError, false),
	}

	for _, test := range tests {
		test.Test(t)
	}

}
