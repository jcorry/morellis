package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jcorry/morellis/pkg/models"

	"github.com/bmizerany/pat"

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

func TestNewPermissionsCheck(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	tests := []struct {
		name       string
		reqUUID    bool
		permission []string
		wantCode   int
	}{
		{"Valid permission", false, []string{"user:read"}, 200},
		{"Invalid permission", false, []string{"foo:bar"}, 401},
		{"Self permission", false, []string{"self:read"}, 200},
		{"Self permission with mismatched UUIDs", true, []string{"self:read"}, 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			handlerToTest := NewPermissionsCheck(nextHandler, tt.permission)

			var reqUUID string
			if tt.reqUUID {
				uid, err = uuid.NewRandom()
				if err != nil {
					t.Fatal(err)
				}
				user.Permissions = []models.UserPermission{
					{
						Permission: models.Permission{Name: "self:read"},
					},
					{
						Permission: models.Permission{Name: "self:write"},
					},
				}

				reqUUID = uid.String()
			} else {
				reqUUID = user.UUID.String()
			}

			reqToken, err := generateToken(user)
			if err != nil {
				t.Fatal(err)
			}

			testUrl, err := url.Parse(fmt.Sprintf("http://testing/testing/%s", reqUUID))

			if err != nil {
				t.Fatal(err)
			}
			req := http.Request{
				Method: http.MethodGet,
				URL:    testUrl,
				Header: map[string][]string{
					"Content-Type":  {"application/json"},
					"Authorization": {fmt.Sprintf("Bearer %s", reqToken)},
				},
			}
			ctx := context.WithValue(req.Context(), ContextKeyUser, user)

			w := httptest.NewRecorder()

			UserRouter(handlerToTest).ServeHTTP(w, req.WithContext(ctx))
			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("Want code %d, Got code %d; User UUID %s; Req UUID %s", tt.wantCode, w.Result().StatusCode, user.UUID.String(), reqUUID)
				t.Errorf("User Permissions: %v", user.Permissions)
			}
		})
	}
}

func UserRouter(handler http.Handler) *pat.PatternServeMux {
	mux := pat.New()
	mux.Get("/testing/:uuid", handler)
	return mux
}
