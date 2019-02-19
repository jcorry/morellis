package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jcorry/morellis/pkg/models"
)

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		app.serverError(w, err)
	}
	defer r.Body.Close()

	user, err = app.users.Insert(user.FirstName, user.LastName, user.Email, user.Phone)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, user)
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	var updateUser *models.User

	if err != nil {
		app.notFound(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&updateUser)
	updateUser.ID = int64(id)

	if err != nil {
		app.serverError(w, err)
	}

	user, err := app.users.Update(updateUser)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, user)
}

// Get a single user by ID.
func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	user, err := app.users.Get(id)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, user)
}

func (app *application) listUser(w http.ResponseWriter, r *http.Request) {
	var err error
	params := r.URL.Query()

	l := params.Get("count")
	limit := 0
	if l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	o := params.Get("start")
	offset := 0
	if o != "" {
		offset, err = strconv.Atoi(o)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	sb := params.Get("sortBy")

	//sd := params.Get("sortDir")

	users, err := app.users.List(limit, offset, sb)

	if err != nil {
		app.serverError(w, err)
	}

	meta := make(map[string]interface{})
	meta["totalRecords"] = app.users.Count()
	meta["count"] = limit
	meta["start"] = offset
	meta["sortBy"] = sb

	response := make(map[string]interface{})
	response["meta"] = meta
	response["items"] = users

	app.jsonResponse(w, response)
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	res, err := app.users.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	if res {
		app.noContentResponse(w)
	}

	app.notFound(w)
}
