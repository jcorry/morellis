package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/jcorry/morellis/pkg/models"
	"github.com/jcorry/morellis/pkg/models/mysql"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication(t *testing.T) *application {
	db := mysql.NewTestDB(t)
	rdb := mysql.NewTestRedis(t)

	return &application{
		errorLog:    log.New(ioutil.Discard, "", 0),
		infoLog:     log.New(ioutil.Discard, "", 0),
		users:       &mysql.UserModel{DB: db, Redis: rdb},
		stores:      &mysql.StoreModel{DB: db},
		flavors:     &mysql.FlavorModel{DB: db},
		ingredients: &mysql.IngredientModel{DB: db},
		mapsApiKey:  os.Getenv("GMAP_API_KEY"),
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

	switch strings.ToLower(method) {
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
			ID:     4,
			Status: "Verified",
			UUID:   uid,
			Permissions: []models.UserPermission{
				{
					17,
					models.Permission{ID: 1, Name: "user:read"},
				},
				{
					24,
					models.Permission{ID: 2, Name: "user:write"},
				},
			},
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

type FakeResponse struct {
	t       *testing.T
	headers http.Header
	body    []byte
	status  int
}

func NewFakeResponse(t *testing.T) *FakeResponse {
	return &FakeResponse{
		t:       t,
		headers: make(http.Header),
	}
}

func (r *FakeResponse) Header() http.Header {
	return r.headers
}

func (r *FakeResponse) Write(body []byte) (int, error) {
	r.body = body
	return len(body), nil
}

func (r *FakeResponse) WriteHeader(status int) {
	r.status = status
}

func (r *FakeResponse) Assert(status int, body string) {
	if r.status != status {
		r.t.Errorf("expected status %+v to equal %+v", r.status, status)
	}
	if string(r.body) != string(body) {
		r.t.Errorf("expected body %v to equal %v", string(r.body), body)
	}
}
