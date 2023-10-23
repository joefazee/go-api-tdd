package main

import (
	"encoding/json"
	"github.com/joefazee/go-api-tdd/pkg/common"
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"github.com/joefazee/go-api-tdd/pkg/security"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateBlog(t *testing.T) {

	defaultJWT, err := security.NewJWT(key)
	if err != nil {
		t.Fatal(err)
	}
	srv := newServer(testStore, defaultJWT)
	ts := newTestServer(srv.routes())

	user, err := srv.store.CreateUser(&domain.User{
		Email:    "test@yy.com",
		Password: "password",
		Name:     "OJ",
	})

	if err != nil {
		t.Fatal(err)
	}

	jwtPayload, _ := srv.jwt.CreateToken(*user, 1*time.Minute)

	testCases := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "Created",
			body:         `{ "title": "Welcome to our first blog","body": "Here we go!!!"}`,
			expectedCode: http.StatusCreated,
		},

		{
			name:         "invalid json",
			body:         `{ invalid"}`,
			expectedCode: http.StatusBadRequest,
		},

		{
			name:         "validation error",
			body:         `{ "title": "","body": ""}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/blog/create", strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer "+jwtPayload.Token)
			req.RequestURI = ""

			res, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if tc.expectedCode != res.StatusCode {
				t.Errorf("want status code of %d; got %d", tc.expectedCode, res.StatusCode)
			}
		})
	}

	_ = srv.store.DeleteUserByID(user.ID)
	_ = srv.store.DeleteAllPosts()
}

func Test_GetAllUserPosts(t *testing.T) {

	defaultJWT, err := security.NewJWT(key)
	if err != nil {
		t.Fatal(err)
	}
	srv := newServer(testStore, defaultJWT)
	ts := newTestServer(srv.routes())

	user, err := srv.store.CreateUser(&domain.User{
		Email:    common.RandomString(10) + "@yy.com",
		Password: "password",
		Name:     "OJ",
	})

	if err != nil {
		t.Fatal(err)
	}

	jwtPayload, _ := srv.jwt.CreateToken(*user, 1*time.Minute)

	for i := 0; i < 2; i++ {
		_, err = srv.store.CreatePost(&domain.Post{
			UserID: user.ID,
			Title:  common.RandomString(20),
			Body:   common.RandomString(20),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, ts.URL+"/api/v1/blog/list", nil)
	req.Header.Set("Authorization", "Bearer "+jwtPayload.Token)
	req.RequestURI = ""

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("want status code of %d; got %d", http.StatusOK, res.StatusCode)
	}

	var postRes = struct {
		Posts []domain.Post `json:"posts"`
	}{}

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err = json.Unmarshal(resBodyBytes, &postRes); err != nil {
		t.Fatal(err)
	}

	if len(postRes.Posts) != 2 {
		t.Errorf("expected 2 posts for this user; got %d", len(postRes.Posts))
	}

}
