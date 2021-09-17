package main

import (
	"net/http"
)

func prepareResponseHeaders(w http.ResponseWriter) {
	for key, val := range defaultResponseHeaders {
		w.Header().Set(key, val)
	}
}
