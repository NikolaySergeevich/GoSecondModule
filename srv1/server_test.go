package main

import (
	"net/http"
	"testing"
)

func TestGetAll(t *testing.T) {
	resp, err := http.Get("http://localhost:8081/get_all")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if resp.Status != "200 OK" {
		t.Fatalf("код ответа равен " + resp.Status)
	}
}
