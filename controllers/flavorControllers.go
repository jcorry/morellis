package controllers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/matryer/respond.v1"

	"github.com/jcorry/morellis/models"
)

var CreateFlavor = func(w http.ResponseWriter, r *http.Request) {
	flavor := &models.Flavor{}
	err := json.NewDecoder(r.Body).Decode(flavor)
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
	}
	respond.With(w, r, http.StatusOK, flavor)
}

var GetFlavors = func(w http.ResponseWriter, r *http.Request) {
	flavors, err := models.GetFlavors()
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
	}
	respond.With(w, r, http.StatusOK, flavors)
}
