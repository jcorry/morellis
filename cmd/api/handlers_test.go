package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/jcorry/morellis/pkg/models/mock"

	"github.com/jcorry/morellis/pkg/models"
)

func TestCreateAuth(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		email    string
		password string
		wantBody []byte
		wantCode int
	}{
		{"Valid credentials", "valid@example.com", "password", []byte(`{"Token":`), 200},
		{"Invalid credentials", "noauth@example.com", "password", []byte("Not Found"), 404},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"email":    tt.email,
				"password": tt.password,
			}

			reqBytes, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}

			code, _, body := ts.request(t, "post", "/api/v1/auth", bytes.NewBuffer(reqBytes))

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}

}

func TestCreateUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name      string
		firstName string
		lastName  string
		phone     string
		email     string
		password  string
		wantCode  int
		wantBody  []byte
	}{
		{"Valid submission", "Bob", "McTestFace", "867-5309", "bob@testy.com", "valid-password", http.StatusOK, []byte("Bob")},
		{"Duplicate email", "Bob", "McTestFace", "867-5309", "dupe@example.com", "valid-password", http.StatusBadRequest, []byte("duplicate email")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"firstName": tt.firstName,
				"lastName":  tt.lastName,
				"phone":     tt.phone,
				"email":     tt.email,
				"password":  tt.password,
			}

			reqBytes, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}

			code, _, body := ts.request(t, "post", "/api/v1/user", bytes.NewBuffer(reqBytes))

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestPartialUpdateUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name      string
		id        int
		firstName string
		lastName  string
		phone     string
		email     string
		password  string
		wantCode  int
		wantBody  []byte
	}{
		{"Valid submission", 1, "Bob", "McTestFace", "867-5309", "bob@testy.com", "valid-password", http.StatusOK, []byte("Bob")},
		{"Duplicate email", 1, "Bob", "McTestFace", "867-5309", "dupe@example.com", "valid-password", http.StatusBadRequest, []byte("duplicate email")},
		{"Invalid ID", 0, "Bob", "McTestFace", "867-5309", "dupe@example.com", "valid-password", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"firstName": tt.firstName,
				"lastName":  tt.lastName,
				"phone":     tt.phone,
				"email":     tt.email,
				"password":  tt.password,
			}

			reqBytes, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}

			urlPath := fmt.Sprintf("/api/v1/user/%d", tt.id)
			code, _, body := ts.request(t, "patch", urlPath, bytes.NewBuffer(reqBytes))

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name     string
		id       string
		wantCode int
		wantBody []byte
	}{
		{"Valid request", id.String(), http.StatusOK, []byte("McTestFace")},
		{"No record", "foo", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/api/v1/user/%s", tt.id)

			code, _, body := ts.get(t, urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestGetStore(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		id       int
		wantCode int
		wantBody []byte
	}{
		{"Valid request", 1, http.StatusOK, []byte("Test Store")},
		{"No record", 0, http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/api/v1/store/%d", tt.id)
			code, _, body := ts.get(t, urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestCreateFlavor(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name              string
		flavorName        string
		flavorDescription string
		wantCode          int
		wantBody          []byte
	}{
		{
			"Valid Flavor",
			"Flava' Flav",
			"A new flavor that is just delicious.",
			200,
			[]byte("Flava' Flav"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":        tt.flavorName,
				"description": tt.flavorDescription,
				"ingredients": []*models.Ingredient{
					{
						Name: "chocolate",
					},
					{
						Name: "sriracha",
					},
				},
			}

			reqBytes, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}

			code, _, body := ts.request(t, "post", "/api/v1/flavor", bytes.NewBuffer(reqBytes))

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestGetFlavor(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		id       int64
		wantCode int
	}{
		{"Valid flavor ID", 1, 200},
		{"Invalid flavor ID", 100, 404},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/api/v1/flavor/%d", tt.id)
			wantBody := ""
			if tt.wantCode == 200 {
				wantBody = fmt.Sprintf(`{"id":%d,`, tt.id)
			} else {
				wantBody = "Not Found"
			}

			code, _, body := ts.get(t, urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, []byte(wantBody)) {
				t.Errorf("want body %s to contain %q", body, wantBody)
			}
		})
	}
}

func TestListFlavor(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	urlPath := "/api/v1/flavor"
	code, _, body := ts.get(t, urlPath)

	if code != 200 {
		t.Errorf("want %d, got %d", 200, code)
	}

	wantString := mock.MockFlavors[1].Name

	if !bytes.Contains(body, []byte(fmt.Sprintf(`"name":"%s"`, wantString))) {
		t.Errorf("want %s, got %s", wantString, body)
	}
}

func TestListStore(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	urlPath := "/api/v1/store"
	code, _, body := ts.get(t, urlPath)

	if code != 200 {
		t.Errorf("want %d, got %d", 200, code)
	}

	wantString := mock.MockStores[0].Name

	if !bytes.Contains(body, []byte(fmt.Sprintf(`"name":"%s"`, wantString))) {
		t.Errorf("want %s, got %s", wantString, body)
	}
}
