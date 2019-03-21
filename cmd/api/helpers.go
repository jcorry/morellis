package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/jcorry/morellis/pkg/models"

	"github.com/dgrijalva/jwt-go"
)

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
	type Claims struct {
		UUID string `json:"uuid"`
		jwt.StandardClaims
	}

	claims := Claims{
		user.UUID.String(),
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

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	key, err = jwt.ParseRSAPrivateKeyFromPEM(pem.EncodeToMemory(privateKey))
	fatal(err)

	return key, nil
}

func getVerifyKey(privateKey *rsa.PublicKey) (*rsa.PublicKey, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)
	fatal(err)

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	key, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	fatal(err)

	return key, nil
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
