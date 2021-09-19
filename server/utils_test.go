package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractStringFromMap(t *testing.T) {

	tests := []struct {
		valName string
		val     string
		err     string
		m       TestMap
	}{
		{"name", "test", "", TestMap{"name": "test"}},
		{"name", "", "", TestMap{"name": ""}},
		{"name", "", ExtractErrInvalidType, TestMap{"name": 52}},
		{"name", "", ExtractErrNotFound, TestMap{"test": 52}},
		{"name", "", ExtractErrNotFound, TestMap{"test": 52}},
	}

	for _, test := range tests {
		val, err := extractStringFromMap(test.valName, test.m)

		if err != nil {
			if err.Error() != test.err {
				t.Errorf("Got error %v expected %v", err.Error(), test.err)
			}
		} else {
			if len(test.err) > 0 {
				t.Errorf("Got error %v expected %v", nil, test.err)
			}
		}

		assert.Equal(t, test.val, val)
	}
}
