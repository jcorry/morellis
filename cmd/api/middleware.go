package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	claims := r.Context().Value(ContextKeyUser).(*Claims)
	reqUUID := r.URL.Query().Get(":uuid")
	if !checkPermissions(claims, pc.permissions, reqUUID) {
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
					ctx := context.WithValue(r.Context(), ContextKeyUser, claims)
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

func checkPermissions(c *Claims, permissions []string, reqUUID string) bool {
	for _, userPermission := range c.Permissions {
		for _, requiredPermission := range permissions {
			if userPermission.Permission.Name == requiredPermission {
				if requiredPermission == "user:read" || requiredPermission == "user:write" {
					return true
				}

				if (requiredPermission == "self:read" || requiredPermission == "self:write") &&
					c.UUID == reqUUID {
					return true
				}
			}
		}
	}
	return false
}
