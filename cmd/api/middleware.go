package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
)

var ContextKeyUser = "AuthUser"

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
					user.Permissions, err = app.users.GetPermissions(int(user.ID))
					if err != nil {
						app.errorLog.Output(2, err.Error())
						w.WriteHeader(http.StatusUnauthorized)
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
