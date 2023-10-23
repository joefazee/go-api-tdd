package main

import (
	"encoding/json"
	"github.com/joefazee/go-api-tdd/pkg/common"
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"github.com/joefazee/go-api-tdd/pkg/security"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {

	testCases := []struct {
		name         string
		expectedCode int
		body         string
	}{
		{
			name:         "OK",
			expectedCode: http.StatusOK,
			body: `{
			"name": "Joseph",
			"email": "aj2@joseh.com",
			"password": "password"}`,
		},

		{
			name:         "Bad JSON",
			expectedCode: http.StatusBadRequest,
			body:         `{"name": "}`,
		},

		{
			name:         "validation error",
			expectedCode: http.StatusBadRequest,
			body: `{
			"name": "",
			"email": "",
			"password": ""}`,
		},
	}

	srv := newServer(testStore, nil)
	ts := newTestServer(srv.routes())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			res, err := ts.Client().Post(ts.URL+"/api/v1/users/create",
				"application/json",
				strings.NewReader(tc.body))

			if err != nil {
				t.Fatal(err)
			}

			if tc.expectedCode != res.StatusCode {
				t.Errorf("want status code of %d; got %d", tc.expectedCode, res.StatusCode)
			}
		})
	}

}

func TestUserLogin(t *testing.T) {

	defaultJWT, err := security.NewJWT(key)
	if err != nil {
		t.Fatal(err)
	}
	srv := newServer(testStore, defaultJWT)
	ts := newTestServer(srv.routes())

	_, err = srv.store.CreateUser(&domain.User{
		Email:    common.RandomString(10) + "@yy.com",
		Password: "password",
		Name:     "OJ",
	})

	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		expectedCode int
		body         string
		checkBody    func(t *testing.T, body []byte)
	}{
		{
			name: "OK",
			body: `{
			"email": "test@yy.com",
			"password": "password"
				}`,
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				loginRes := struct {
					Token string `json:"token"`
				}{}

				if err := json.Unmarshal(body, &loginRes); err != nil {
					t.Fatal(err)
				}
				if loginRes.Token == "" {
					t.Errorf("want token; got empty string in loginRes.Token")
				}
			},
		},
		{
			name:         "Invalid json",
			body:         `{X"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validation error",
			body: `{
			"email": "",
			"password": ""
				}`,
			expectedCode: http.StatusBadRequest,
		},

		{
			name: "user not found",
			body: `{
			"email": "invalid@y.com",
			"password": "sssss"
				}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid password",
			body: `{
			"email": "test@yy.com",
			"password": "invalid password"
				}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ts.Client().
				Post(ts.URL+"/api/v1/users/login",
					"application/json",
					strings.NewReader(tc.body),
				)
			if err != nil {
				t.Fatal(err)
			}

			if tc.expectedCode != res.StatusCode {
				t.Errorf("expect status of %d; got %d", tc.expectedCode, res.StatusCode)
			}

			if tc.checkBody != nil {
				defer res.Body.Close()

				bodyBS, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				tc.checkBody(t, bodyBS)
			}
		})
	}

}
