package controllers

import (
	"encoding/json"
	"morellis/models"
	u "morellis/utils"
	"net/http"

	"github.com/jinzhu/gorm"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	accountStatus := &models.AccountStatus{}
	err := models.GetDB().Table("account_statuses").Where("value = ?", "Pending").First(accountStatus).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Message(false, "Email address not found")
			return
		}
		u.Message(false, "Connection error. Please retry")
		return
	}

	account := &models.Account{AccountStatus: *accountStatus}

	err = json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Unable to parse account."))
		return
	}

	resp := account.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}
