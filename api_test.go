package main

import (
	"bytes"
	"context"
	"encoding/json"
	api "fenix/api"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
)

const endpoint = "http://localhost:8080"

func setupTestCase(t *testing.T) func(bool) {
	testID, _ := uuid.NewUUID()

	password, ok := os.LookupEnv("FENIX_PASS")

	if !ok {
		t.Fatal("No FENIX_PASS env var set!")
	}

	a := api.NewAPI("testing", password, testID.String())

	a.Serve(true)

	return func(ok bool) {
		cli, _ := a.UserDatabase.Database()

		cli.Delete(context.Background(), "/testing/"+testID.String()+"*")

		cli.Close()

		if ok {
			os.Exit(0)
		}
		os.Exit(1)
	}
}

func CreateUser(bodyJSON, email, password string) (int, map[string]interface{}, error) {
	b := bytes.Buffer{}

	b.Write([]byte(bodyJSON))

	req, err := http.NewRequest("POST", endpoint+"/6.0.1/create", &b)

	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(email, password)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return 0, nil, err
	}

	var body map[string]interface{}

	bodyBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return 0, nil, err
	}

	err = json.Unmarshal(bodyBytes, &body)

	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, body, nil
}

type createCase struct {
	name, email, password, json string
	okResult                    bool
}

func newCreateCase(json, email, password, name string, okResult bool) createCase {
	c := createCase{}
	c.name = name
	c.json = json
	c.email = email
	c.password = password
	c.okResult = okResult

	return c
}

func TestCreate(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	var testCases [4]createCase

	testCases[0] = newCreateCase("{\"username\":\"Rick Astley\"}", "rick.astley@rickroll.com", "yay", "CreateUser", true)
	testCases[1] = newCreateCase("{\"username\":\"Rick Astley\"}", "rick.astley@rickroll.com", "yay", "CreateDuplicateUser", false)
	testCases[2] = newCreateCase("{\"username\":\"Rick Astley Likes To Rickroll People With Long Usernames That Are Over 32 Characters\"}", "astley.rick@rickroll.com", "yay", "CreateUserWithLongName", false)
	testCases[3] = newCreateCase("invalid json", "invalid.astley@rickroll.com", "yay", "CreateUserWithInvalidJson", false)

	for testIndex, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			statusCode, body, err := CreateUser(testCase.json, testCase.email, testCase.password)
			// This should never return an error, since any errors we expect are just API errors.  If this happens, something is wrong.
			if err != nil {
				// Values are surrounded by zero width spaces (​) (Unicode U+200B)
				t.Errorf("Failed on ​%v​.  Status code: ​%v​.  Body: ​%v​.  Error: ​%v​", testIndex, statusCode, body, err)

				return
			}

			// 201 Creates is the correct response for this.
			if (statusCode == 201) == testCase.okResult {
				t.Logf("%v ok", testIndex)
			} else {
				expected := "succeed"
				if !testCase.okResult {
					expected = "fail"
				}

				t.Logf("Failed on ​%v​.  Status code: ​%v​.  Body: ​%v​.  Expected function to %v", testIndex, statusCode, body, expected)
				t.Fail()
			}
		})
	}

	teardownTestCase(!t.Failed())
}
