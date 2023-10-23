package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"github.com/joefazee/go-api-tdd/pkg/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApplyAuthentication(t *testing.T) {

	defaultJWT, err := security.NewJWT(key)
	if err != nil {
		t.Fatal(err)
	}
	srv := newServer(testStore, defaultJWT)
	ts := newTestServer(srv.routes())

	user, err := srv.store.CreateUser(&domain.User{
		Email:    "test2@yy.com",
		Password: "password",
		Name:     "OJ",
	})

	testCases := []struct {
		name         string
		expectedCode int
		setupHeader  func(r *http.Request)
		checkBody    func(t *testing.T, res *http.Response)
	}{
		{
			name:         "OK",
			expectedCode: http.StatusOK,
			setupHeader: func(r *http.Request) {
				jwtPayload, _ := srv.jwt.CreateToken(*user, 1*time.Minute)
				r.Header.Set("Authorization", "Bearer "+jwtPayload.Token)

			},
		},
		{
			name:         "No auth header",
			expectedCode: http.StatusUnauthorized,
			setupHeader: func(r *http.Request) {
			},
		},
		{
			name:         "No auth header",
			expectedCode: http.StatusUnauthorized,
			setupHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "invalid")
			},
		},

		{
			name:         "Expired token",
			expectedCode: http.StatusUnauthorized,
			setupHeader: func(r *http.Request) {
				jwtPayload, _ := srv.jwt.CreateToken(*user, -1*time.Minute)
				r.Header.Set("Authorization", "Bearer "+jwtPayload.Token)

			},
		},

		{
			name:         "invalid token",
			expectedCode: http.StatusUnauthorized,
			setupHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid")
			},
		},

		{
			name:         "invalid user id",
			expectedCode: http.StatusUnauthorized,
			setupHeader: func(r *http.Request) {
				jwtPayload, _ := srv.jwt.CreateToken(domain.User{
					ID:    -122,
					Name:  "test",
					Email: "test@tst.com",
				}, 1*time.Minute)
				r.Header.Set("Authorization", "Bearer "+jwtPayload.Token)
			},
		},
	}

	srv.router.GET("/auth", srv.applyAuthentication(), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{})
	})

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ts.URL+"/auth", nil)
			req.RequestURI = ""
			tc.setupHeader(req)

			res, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if tc.expectedCode != res.StatusCode {
				t.Errorf("want %d; got %d", tc.expectedCode, res.StatusCode)
			}

			if tc.checkBody != nil {
				tc.checkBody(t, res)
			}

		})
	}
}
