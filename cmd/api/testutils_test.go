package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/jcorry/morellis/pkg/models"
	"github.com/jcorry/morellis/pkg/models/mock"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(ioutil.Discard, "", 0),
		infoLog:  log.New(ioutil.Discard, "", 0),
		users:    &mock.UserModel{},
		stores:   &mock.StoreModel{},
		flavors:  &mock.FlavorModel{},
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	return &testServer{ts}
}

func (ts *testServer) request(t *testing.T, method string, urlPath string, reqBody io.Reader, authorized bool) (int, http.Header, []byte) {
	u, err := url.Parse(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	req := &http.Request{
		Method: "PATCH",
		URL:    u,
		Body:   ioutil.NopCloser(reqBody),
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	switch method {
	case "get":
		req.Method = "GET"
	case "post":
		req.Method = "POST"
	case "patch":
		req.Method = "PATCH"
	case "put":
		req.Method = "PUT"
	case "delete":
		req.Method = "DELETE"
	}

	if authorized {
		uid, err := uuid.NewRandom()
		if err != nil {
			t.Fatal(err)
		}

		user := models.User{
			UUID: uid,
		}

		token, err := generateToken(&user)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
