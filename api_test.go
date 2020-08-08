package main

import (
	"bytes"
	"encoding/json"
	api "fenix/api"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

const endpoint = "http://localhost:8080"

func setupTestCase(t *testing.T) func(t *testing.T) {

	password, ok := os.LookupEnv("FENIX_PASS")

	if !ok {
		t.Fatal("No FENIX_PASS env var set!")
	}

	a := api.NewAPI("testing", password)
	err := make(chan api.Error)

	go a.Serve(err, true)

	go func() {
		for {
			t.Error(<-err)
		}
	}()

	return func(t *testing.T) {
		os.Exit(0)
	}
}

// TestCreate tests /v6.0.1/create with a valid request
func TestValidCreate(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	b := bytes.Buffer{}
	bytes, err := json.Marshal(map[string]string{"username": "Rick Astley"})

	if err != nil {
		t.Error(err)
	}

	b.Write(bytes)

	req, err := http.NewRequest("POST", endpoint+"/6.0.1/create", &b)

	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("rick.astley@rickroll.com", "ilikerickrolls")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	var body *map[string]interface{}

	bodyBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Error(err)
	}

	json.Unmarshal(bodyBytes, body)

	if res.StatusCode != 201 {
		t.Logf("FAIL: Status code: %v", res.StatusCode)
		t.Logf("FAIL: Body: %v", body)
		t.FailNow()
	}
}
