package main

import (
	"math/rand"
	"net/http"
)

func prepareResponseHeaders(w http.ResponseWriter) {
	for key, val := range defaultResponseHeaders {
		w.Header().Set(key, val)
	}
}

func extractStringFromMap(val string, m map[string]interface{}) (s string, err error) {
	valI, ok := m[val]
	if !ok {
		return
	}

	s, ok = valI.(string)
	if !ok {
		return
	}

	return
}

func randRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}
