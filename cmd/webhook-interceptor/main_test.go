package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(health))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"health":"ok"}`
	if string(body) != want {
		t.Fatalf("Want %s, got %s", want, string(body))
	}
}
