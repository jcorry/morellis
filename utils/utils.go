package utils

import (
	"encoding/json"
	"net/http"

	"gopkg.in/matryer/respond.v1"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RespondWithObject(w http.ResponseWriter, r *http.Request, t string, obj interface{}) {
	respond.With(w, r, http.StatusOK, obj)
}

func RespondWithCollection(w http.ResponseWriter, r *http.Request, t string, collection []interface{}) {
	meta := map[string]interface{}{}
	meta["count"] = len(collection)

	resp := map[string]interface{}{}
	resp["meta"] = meta
	resp["items"] = collection

	respond.With(w, r, http.StatusOK, resp)
}
