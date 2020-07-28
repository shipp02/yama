package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	var reqData = struct {
		Name         string
		Username     string
		PasswordHash string
	}{
		Name:         "new-name",
		Username:     "test=user",
		PasswordHash: "Very good pass",
	}

	binReq, _ := json.Marshal(reqData)
	req := bytes.NewBuffer(binReq)
	response, err := http.Post("http://localhost:8080/edit/u/create",
		"application/json", req)
	contentLength := response.ContentLength
	responseBytes := make([]byte, contentLength)
	_, err = response.Body.Read(responseBytes)
	if err != nil {
		log.Println(err)
		t.Fail()
		return
	}
	resp := string(responseBytes)
	if strings.Contains(resp, "User exists") {
		log.Println("user exists")
		t.Fail()
	}

	response, err = http.Post("http://localhost:8080/edit/u/create",
		"application/json", req)
	contentLength = response.ContentLength
	responseBytes = make([]byte, contentLength)
	_, err = response.Body.Read(responseBytes)
	if err != nil {
		log.Println(err)
		t.Fail()
		return
	}
	resp = string(responseBytes)
	if !strings.Contains(resp, "User exists") {
		log.Println("User not created in first time")
		t.Fail()
	}
}
