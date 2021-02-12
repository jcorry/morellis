package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jcorry/morellis/pkg/models"
)

func TestSmsAuthRequest(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		validSignature bool
		statusCode     int
	}{
		{
			validSignature: true,
			statusCode:     http.StatusOK,
		},
		{
			validSignature: false,
			statusCode:     http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			code, _, _ := ts.request(t, "post", "/webhooks/v1/sms/auth", bytes.NewBuffer([]byte{}), false)
			if code != tt.statusCode {
				t.Errorf("unexpected code; want %d; got %d", tt.statusCode, code)
			}
		})
	}

}

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
		{"Valid credentials should have token", "alice@example.com", "password", []byte(`{"token":`), 200},
		{"Valid credentials should expire", "alice@example.com", "password", []byte(`"expires":`), 200},
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

			code, _, body := ts.request(t, "post", "/api/v1/auth", bytes.NewBuffer(reqBytes), false)

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
		name        string
		firstName   string
		lastName    string
		phone       string
		email       string
		password    string
		permissions []string
		wantCode    int
		wantBody    []byte
	}{
		{"Valid submission", "Bob", "McTestFace", "867-5310", "bob@testy.com", "valid-password", []string{"user:read", "user:write"}, http.StatusOK, []byte("Bob")},
		{"Duplicate email", "Bob", "McTestFace", "867-5311", "bob@testy.com", "valid-password", []string{"user:read", "user:write"}, http.StatusBadRequest, []byte("Duplicate email")},
		{"Duplicate phone", "Bob", "McTestFace", "867-5309", "dupe@example.com", "valid-password", []string{"user:read", "user:write"}, http.StatusBadRequest, []byte("Duplicate phone")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			permissions := []models.UserPermission{}

			for _, n := range tt.permissions {
				p := models.Permission{Name: n}
				up := models.UserPermission{Permission: p}
				permissions = append(permissions, up)
			}

			reqBody := map[string]interface{}{
				"firstName":   tt.firstName,
				"lastName":    tt.lastName,
				"phone":       tt.phone,
				"email":       tt.email,
				"password":    tt.password,
				"permissions": permissions,
				"status":      "verified",
			}

			reqBytes, err := json.Marshal(reqBody)
			t.Log(fmt.Sprintf("%s", reqBytes))
			if err != nil {
				t.Fatal(err)
			}

			code, _, body := ts.request(t, "post", "/api/v1/user", bytes.NewBuffer(reqBytes), true)

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
			code, _, body := ts.request(t, "patch", urlPath, bytes.NewBuffer(reqBytes), true)

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

	// get a valid userID
	u, err := app.users.GetByPhone("867-5309")
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name     string
		id       string
		wantCode int
		wantBody []byte
	}{
		{"Valid request", u.UUID.String(), http.StatusOK, []byte("alice@example.com")},
		{"Invalid UUID", "foo", http.StatusNotFound, nil},
		{"No record", "e6fc6b5a-882c-40ba-b860-b11a413ec2df", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/api/v1/user/%s", tt.id)

			code, _, body := ts.request(t, "get", urlPath, bytes.NewBuffer(nil), true)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// get a valid userID
	u, err := app.users.GetByPhone("867-5309")
	if err != nil {
		t.Error(err)
	}

	urlPath := fmt.Sprintf("/api/v1/user/%s", u.UUID.String())

	code, _, _ := ts.request(t, "delete", urlPath, bytes.NewBuffer(nil), true)

	if code != http.StatusNoContent {
		t.Errorf("want %d; got %d", http.StatusNoContent, code)
	}
}

func TestHandlers_CreateStore(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		store    models.Store
		wantCode int
		wantBody []byte
	}{
		{
			"Valid Store",
			models.Store{
				Name:    "Foo",
				Phone:   "867-5309",
				Email:   "foo@bar.com",
				URL:     "http://www.bar.com",
				Address: "1600 Pennsylvania Ave",
				City:    "Washington",
				State:   "DC",
				Zip:     "20500-0004",
			},
			200,
			[]byte("Pennsylvania"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.store)
			if err != nil {
				t.Errorf("Unexpected err; %s", err)
			}

			code, _, body := ts.request(t, "post", "/api/v1/store", bytes.NewBuffer(reqBody), true)

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
		{"Valid request", 1, http.StatusOK, []byte("Moreland")},
		{"No record", 0, http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/api/v1/store/%d", tt.id)
			code, _, body := ts.request(t, "get", urlPath, bytes.NewBuffer(nil), true)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestActivateStoreFlavor(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name      string
		qStoreID  int
		qFlavorID int
		bStoreID  int
		bFlavorID int
		position  int
		wantCode  int
		wantBody  []byte
	}{
		{"Successful Activation", 1, 2, 1, 2, 3, 200, []byte("")},
		{"Mismatched Store IDs", 1, 1, 2, 1, 3, http.StatusBadRequest, []byte("")},
		{"Mismatched Flavor IDs", 1, 1, 1, 2, 3, http.StatusBadRequest, []byte("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/api/v1/store/%d/flavor/%d", tt.qStoreID, tt.qFlavorID)
			reqBody := map[string]interface{}{
				"position":  tt.position,
				"store_id":  tt.bStoreID,
				"flavor_id": tt.bFlavorID,
			}
			reqBytes, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}
			code, _, body := ts.request(t, "post", urlPath, bytes.NewBuffer(reqBytes), true)

			if code != tt.wantCode {
				t.Errorf("Want %d; Got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestDeactivateStoreFlavor(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		storeID  int
		flavorID int
		wantCode int
	}{
		{"Successful deactivation", 1, 1, 204},
		{"Deactivation: no rows affected", 3, 3, 204},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/api/v1/store/%d/flavor/%d", tt.storeID, tt.flavorID)
			code, _, _ := ts.request(t, "delete", urlPath, bytes.NewBuffer(nil), true)

			if code != tt.wantCode {
				t.Errorf("Want %d; Got %d", tt.wantCode, code)
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

			code, _, body := ts.request(t, "post", "/api/v1/flavor", bytes.NewBuffer(reqBytes), true)

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

			code, _, body := ts.request(t, "get", urlPath, bytes.NewBuffer(nil), true)

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
	code, _, _ := ts.request(t, "get", urlPath, bytes.NewBuffer(nil), true)

	if code != 200 {
		t.Errorf("want %d, got %d", 200, code)
	}
}

func TestListStore(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	urlPath := "/api/v1/store"
	code, _, _ := ts.request(t, "get", urlPath, bytes.NewBuffer(nil), true)

	if code != 200 {
		t.Errorf("want %d, got %d", 200, code)
	}
}

func TestUserIngredientAssociations(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	u, err := app.users.GetByPhone("867-5309")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	url := fmt.Sprintf("/api/v1/user/%s/ingredient", u.UUID.String())

	req := map[string]interface{}{
		"userUuid":     u.UUID.String(),
		"ingredientId": 2,
		"storeId":      1,
		"keyword":      "coconut",
	}

	t.Run("create an association", func(t *testing.T) {
		reqBytes, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}

		code, _, body := ts.request(t, "post", url, bytes.NewBuffer(reqBytes), true)
		if code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}

		if !bytes.Contains(body, []byte("coconut")) {
			t.Errorf("want body %s to contain %q", body, []byte("coconut"))
		}
	})

	t.Run("list associations", func(t *testing.T) {
		code, _, body := ts.request(t, http.MethodGet, url, bytes.NewBuffer([]byte{}), true)
		if code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}
		if !bytes.Contains(body, []byte(`"meta":{"count":1,"totalRecords":1}`)) {
			t.Errorf("unexpected body: got %s", body)
		}
	})

	t.Run("create another association", func(t *testing.T) {
		req["ingredientId"] = 1
		reqBytes, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}

		code, _, body := ts.request(t, "post", url, bytes.NewBuffer(reqBytes), true)
		if code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}

		if !bytes.Contains(body, []byte("coconut")) {
			t.Errorf("want body %s to contain %q", body, []byte("coconut"))
		}
	})

	t.Run("list associations", func(t *testing.T) {
		code, _, body := ts.request(t, http.MethodGet, url, bytes.NewBuffer([]byte{}), true)
		if code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}
		if !bytes.Contains(body, []byte(`"meta":{"count":2,"totalRecords":2}`)) {
			t.Errorf("unexpected body: got %s", body)
		}
	})

	t.Run("create an invalid association", func(t *testing.T) {
		req["ingredientId"] = 42
		reqBytes, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}

		code, _, _ := ts.request(t, "post", url, bytes.NewBuffer(reqBytes), true)
		if code != http.StatusNotFound {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}
	})

	t.Run("delete an association", func(t *testing.T) {
		url = fmt.Sprintf("/api/v1/user/%s/ingredient/%d", u.UUID.String(), 1)
		code, _, _ := ts.request(t, "delete", url, bytes.NewBuffer([]byte{}), true)
		if code != http.StatusNoContent {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}
	})

	t.Run("list associations", func(t *testing.T) {
		url = fmt.Sprintf("/api/v1/user/%s/ingredient", u.UUID.String())
		code, _, body := ts.request(t, http.MethodGet, url, bytes.NewBuffer([]byte{}), true)
		if code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}
		if !bytes.Contains(body, []byte(`"meta":{"count":1,"totalRecords":1}`)) {
			t.Errorf("unexpected body: got %s", body)
		}
	})
}
