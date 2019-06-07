package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jcorry/morellis/pkg/models"

	"github.com/dgrijalva/jwt-go"
)

var ContextKeyUser = "AuthUser"

type PermissionsCheck struct {
	handler     http.Handler
	permissions []string
}

func NewPermissionsCheck(handler http.Handler, permissions []string) *PermissionsCheck {
	return &PermissionsCheck{handler, permissions}
}

func (pc *PermissionsCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ContextKeyUser).(*models.User)
	reqUUID := r.URL.Query().Get(":uuid")
	if !checkPermissions(user, pc.permissions, reqUUID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pc.handler.ServeHTTP(w, r)
}

func (app *application) jwtVerification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) == 2 {
				claims := &Claims{}
				token, err := jwt.ParseWithClaims(bearerToken[1], claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
						return nil, fmt.Errorf("No RSA signing method found")
					}

					return verifyKey, nil
				})

				if err != nil {
					app.errorLog.Output(2, err.Error())
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if token.Valid {
					// Is it a valid user?
					uid, err := uuid.Parse(claims.UUID)
					if err != nil {
						app.errorLog.Output(2, err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					user, err := app.users.GetByUUID(uid)
					if err != nil {
						app.errorLog.Output(2, err.Error())
						w.WriteHeader(http.StatusUnauthorized)
						return
					}

					ctx := context.WithValue(r.Context(), ContextKeyUser, user)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func checkPermissions(user *models.User, permissions []string, reqUUID string) bool {
	ok := false
	for _, userPermission := range user.Permissions {
		for _, requiredPermission := range permissions {
			if userPermission.Permission.Name == requiredPermission {
				if requiredPermission == "self:read" || requiredPermission == "self:write" {
					if user.UUID.String() == reqUUID {
						ok = true
					}
				}

				if requiredPermission == "user:read" || requiredPermission == "user:write" {
					ok = true
				}
				return ok
			}
		}
	}
	return false
}
