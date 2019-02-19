package controllers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/matryer/respond.v1"

	"github.com/jcorry/morellis/models"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {

		return
	}
	account, err = account.Create() //Create account
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
	}

	respond.With(w, r, http.StatusCreated, account)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		return
	}

	account, err = models.Login(account.Email, account.Password)
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
	}

	respond.With(w, r, http.StatusOK, account)
}
