package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestMap = map[string]interface{}

type RestTest struct {
	Path               string
	Method             string
	Body, ResponseBody map[string]interface{}
	ResponseCode       int
	IgnoreResponseBody bool
}

func NewRestTest(path, method string, body, responseBody map[string]interface{}, responseCode int, ignoreResponseBody bool) RestTest {
	return RestTest{Path: path, Method: method, Body: body, ResponseBody: responseBody, ResponseCode: responseCode, IgnoreResponseBody: ignoreResponseBody}
}

func (r RestTest) Test(t *testing.T) {

	reqBody, err := json.Marshal(r.Body)
	if err != nil {
		t.Error(err)
	}

	client := http.Client{}

	request, err := http.NewRequest(r.Method, fmt.Sprintf("%v%v", "http://localhost:8000", r.Path), bytes.NewBuffer(reqBody))
	if err != nil {
		t.Error(err)
	}

	res, err := client.Do(request)

	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if !r.IgnoreResponseBody {
		bodyBuff, _ := ioutil.ReadAll(res.Body)
		respBody, err := json.MarshalIndent(r.ResponseBody, "", " ")
		if err != nil {
			t.Error(err)
		}

		temp := map[string]interface{}{}
		err = json.Unmarshal(bodyBuff, &temp)
		if err != nil {
			log.Println(string(bodyBuff))

			t.Error(err)
		}

		delete(temp, "detailed_error")

		fResp, err := json.MarshalIndent(temp, "", " ")
		if err != nil {
			t.Error(err)
		}

		assert.JSONEq(t, string(respBody), string(fResp))
	}
	assert.Equal(t, r.ResponseCode, res.StatusCode)

}

var s Server

func TestMain(m *testing.M) {
	s.Init("./test.db")
	s.Run(false)

	if _, err := s.createUser("Michal"); err != nil {
		log.Fatal(err)
	}

	if _, err := s.createUser("Ania"); err != nil {
		log.Fatal(err)
	}

	if _, err := s.createUser("Dominika"); err != nil {
		log.Fatal(err)
	}

	if _, err := s.createRecord("R1", "t1"); err != nil {
		log.Fatal(err)
	}

	if _, err := s.createRecord("R2", "t1"); err != nil {
		log.Fatal(err)
	}

	if _, err := s.createRecord("R3", "t2"); err != nil {
		log.Fatal(err)
	}

	if err := s.assignRecordToUser(3, 1); err != nil {
		log.Fatal(err)
	}
	if err := s.assignRecordToUser(3, 2); err != nil {
		log.Fatal(err)
	}
	if err := s.assignRecordToUser(3, 3); err != nil {
		log.Fatal(err)
	}

	//	testTableSchemes()

	code := m.Run()
	os.Remove("./test.db")
	os.Exit(code)
}
