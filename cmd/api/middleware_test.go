package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
)

func TestJwtVerificationMiddleware(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name       string
		url        string
		validToken bool
		wantCode   int
	}{
		{"Valid token", "/api/v1/user", true, 200},
		{"No valid token", "/api/v1/user", false, 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u bytes.Buffer
			u.WriteString(string(ts.URL))
			u.WriteString(tt.url)

			url, err := url.Parse(u.String())
			if err != nil {
				t.Fatal(err)
			}

			req := &http.Request{
				Method: "GET",
				URL:    url,
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			}

			if tt.validToken {
				uid, err := uuid.NewRandom()
				if err != nil {
					t.Fatal(err)
				}
				user, err := app.users.GetByUUID(uid)
				if err != nil {
					t.Fatal(err)
				}
				token, err := generateToken(user)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			res, err := ts.Client().Do(req)

			if err != nil {
				t.Fatal(err)
			}
			if res != nil {
				defer res.Body.Close()
			}
			if res.StatusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, res.StatusCode)
			}
		})
	}
}

func TestJwtVerificationAddsUserToContext(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	uid, err := uuid.NewRandom()
	if err != nil {
		t.Fatal(err)
	}

	user, err := app.users.GetByUUID(uid)
	if err != nil {
		t.Fatal(err)
	}

	reqToken, err := generateToken(user)
	if err != nil {
		t.Fatal(err)
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(ContextKeyUser)
		if val == nil {
			t.Error("user not present")
		}
		if val != user {
			t.Error("Not the same user")
		}
	})

	handlerToTest := app.jwtVerification(nextHandler)

	url, err := url.Parse("http://testing")
	if err != nil {
		t.Fatal(err)
	}
	req := http.Request{
		URL: url,
		Header: map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {fmt.Sprintf("Bearer %s", reqToken)},
		},
	}

	handlerToTest.ServeHTTP(httptest.NewRecorder(), &req)
}
