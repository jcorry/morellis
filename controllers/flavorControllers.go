package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jcorry/morellis/models"
	u "github.com/jcorry/morellis/utils"
)

var CreateFlavor = func(w http.ResponseWriter, r *http.Request) {
	flavor := &models.Flavor{}

	err := json.NewDecoder(r.Body).Decode(flavor)
	if err != nil {
		u.Respond(w, u.Message(false, "Error unmarshaling request body"))
		return
	}
	res := flavor.Create()
	u.Respond(w, res)
}

var GetFlavors = func(w http.ResponseWriter, r *http.Request) {
	flavors := models.GetFlavors()
	resp := u.Message(true, "success")
	resp["flavors"] = flavors
	u.Respond(w, resp)
}
