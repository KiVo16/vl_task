package main

import (
	"errors"
	"math/rand"
	"net/http"
)

func prepareResponseHeaders(w http.ResponseWriter) {
	for key, val := range defaultResponseHeaders {
		w.Header().Set(key, val)
	}
}

const (
	ExtractErrNotFound    = "0"
	ExtractErrInvalidType = "1"
)

func extractStringFromMap(val string, m map[string]interface{}) (s string, err error) {
	valI, ok := m[val]
	if !ok {
		return "", errors.New(ExtractErrNotFound)
	}

	s, ok = valI.(string)
	if !ok {
		return "", errors.New(ExtractErrInvalidType)
	}

	return
}

func handleExtractStringFromMapError(w http.ResponseWriter, valName string, err error) {
	switch err.Error() {
	case ExtractErrNotFound:
		NewPredefinedServerError(http.StatusBadRequest, ErrValueNotFound).WithRefersTo(valName).Write(w)
	case ExtractErrInvalidType:
		NewPredefinedServerError(http.StatusBadRequest, ErrValueInvalidType).WithMessage("Expected string").WithRefersTo(valName).Write(w)
	}
}

func randRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}
