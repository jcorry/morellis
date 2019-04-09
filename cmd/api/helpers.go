package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/jcorry/morellis/pkg/models"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UUID        string              `json:"uuid"`
	Permissions []models.Permission `json:"permissions"`
	jwt.StandardClaims
}

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	privKey, err := getSignKey()
	fatal(err)

	signKey = privKey
	verifyKey = &privKey.PublicKey
}

func generateToken(user *models.User) (string, error) {
	claims := Claims{
		user.UUID.String(),
		user.Permissions,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			Issuer:    "morellisicecream.com",
		},
	}

	t := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)

	tokenString, err := t.SignedString(signKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (*Claims, error) {
	c := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("No RSA signing method found")
		}
		return verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return c, err
	}

	return nil, fmt.Errorf("Unable to verify token")
}

func (app *application) badRequest(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(
		w,
		fmt.Sprintf("%s : %s", http.StatusText(http.StatusBadRequest), err.Error()),
		http.StatusBadRequest)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(
		w,
		fmt.Sprintf("%s : %s", http.StatusText(http.StatusInternalServerError), err.Error()),
		http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) jsonResponse(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (app *application) noContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func getSignKey() (*rsa.PrivateKey, error) {
	if signKey != nil {
		return signKey, nil
	}

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	fatal(err)

	err = key.Validate()
	fatal(err)

	return key, nil
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
