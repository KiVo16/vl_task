package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHandlePostUsers(t *testing.T) {

	tests := []RestTest{
		NewRestTest("/users", "POST", TestMap{"name": "test", "type": "t1"}, TestMap{"status": 0, "response": TestMap{"created_id": 4}}, http.StatusOK, false),
		NewRestTest("/users", "POST", TestMap{"name": "test", "type": "t1"}, TestMap{"status": 0, "response": TestMap{"created_id": 5}}, http.StatusOK, false),
		NewRestTest("/users", "POST", TestMap{}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueNotFound,
			"message":    errMap[ErrValueNotFound],
			"refers_to":  "name",
		}, http.StatusBadRequest, false),
		NewRestTest("/users", "POST", TestMap{"name": 25}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    "Expected string",
			"refers_to":  "name",
		}, http.StatusBadRequest, false),
		NewRestTest("/users", "POST", TestMap{"name": true}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    "Expected string",
			"refers_to":  "name",
		}, http.StatusBadRequest, false),
	}

	for _, test := range tests {
		test.Test(t)
	}

}

func TestHandleGetUsers(t *testing.T) {

	tests := []RestTest{
		NewRestTest("/users?limit=1", "GET", TestMap{}, TestMap{"status": 0, "response": []TestMap{TestMap{"id": 1, "name": "Michal"}}}, http.StatusOK, false),
		NewRestTest("/users?limit=1&offset=1", "GET", TestMap{}, TestMap{"status": 0, "response": []TestMap{TestMap{"id": 2, "name": "Ania"}}}, http.StatusOK, false),
		NewRestTest("/users?limit=2&offset=1", "GET", TestMap{}, TestMap{"status": 0, "response": []TestMap{TestMap{"id": 2, "name": "Ania"}, TestMap{"id": 3, "name": "Dominika"}}}, http.StatusOK, false),
		NewRestTest("/users?limit=sd", "GET", TestMap{}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    fmt.Sprintf(errMap[ErrValueInvalidType], "int"),
			"refers_to":  "limit",
		}, http.StatusBadRequest, false),
		NewRestTest("/users?offset=sd", "GET", TestMap{}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    fmt.Sprintf(errMap[ErrValueInvalidType], "int"),
			"refers_to":  "offset",
		}, http.StatusBadRequest, false),
		NewRestTest("/users?limit=1&offset=sd", "GET", TestMap{}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    fmt.Sprintf(errMap[ErrValueInvalidType], "int"),
			"refers_to":  "offset",
		}, http.StatusBadRequest, false),
		NewRestTest("/users?limit=1&offset=true", "GET", TestMap{}, TestMap{
			"http_code":  http.StatusBadRequest,
			"error_code": ErrValueInvalidType,
			"message":    fmt.Sprintf(errMap[ErrValueInvalidType], "int"),
			"refers_to":  "offset",
		}, http.StatusBadRequest, false),
	}

	for _, test := range tests {
		test.Test(t)
	}

}
