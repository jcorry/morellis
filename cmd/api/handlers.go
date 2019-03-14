package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jcorry/morellis/pkg/models"
)

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		app.serverError(w, err)
	}
	defer r.Body.Close()

	user, err = app.users.Insert(user.FirstName, user.LastName, user.Email, user.Phone, user.Password)
	if err != nil {
		if err == models.ErrDuplicateEmail {
			app.badRequest(w, err)
			return
		}

		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, user)
}

func (app *application) partialUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	user, err := app.users.Get(id)

	if err != nil {
		app.notFound(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&user)
	user.ID = int64(id)

	if err != nil {
		app.serverError(w, err)
	}

	user, err = app.users.Update(user)
	if err != nil {
		if err == models.ErrDuplicateEmail {
			app.badRequest(w, err)
			return
		}

		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, user)
}

// Get a single user by ID.
func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromBytes([]byte(r.URL.Query().Get(":uuid")))
	if err != nil || id.String() == "" {
		app.notFound(w)
		return
	}

	user, err := app.users.GetByUUID(id)

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

// Store handlers
func (app *application) createStore(w http.ResponseWriter, r *http.Request) {
	var store *models.Store
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		app.serverError(w, err)
	}

	err = json.Unmarshal(b, &store)

	if err != nil {
		app.serverError(w, err)
	}

	// Geocode the store
	err = app.geocodeStore(store)
	if err != nil {
		app.errorLog.Output(3, err.Error())
	}

	store, err = app.stores.Insert(store.Name, store.Phone, store.Email, store.URL, store.Address, store.City, store.State, store.Zip, store.Lat, store.Lng)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, store)
}

func (app *application) partialUpdateStore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	store, err := app.stores.Get(id)

	if err != nil {
		app.notFound(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&store)

	app.geocodeStore(store)

	store, err = app.stores.Update(id, store.Name, store.Phone, store.Email, store.URL, store.Address, store.City, store.State, store.Zip, store.Lat, store.Lng)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, store)
}

func (app *application) updateStore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	var store *models.Store

	if err != nil {
		app.notFound(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&store)

	app.geocodeStore(store)

	store, err = app.stores.Update(id, store.Name, store.Phone, store.Email, store.URL, store.Address, store.City, store.State, store.Zip, store.Lat, store.Lng)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, store)
}

func (app *application) listStore(w http.ResponseWriter, r *http.Request) {
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

	sb := "s.name"

	stores, err := app.stores.List(limit, offset, sb)
	if err != nil {
		app.serverError(w, err)
	}

	meta := make(map[string]interface{})
	meta["totalRecords"] = app.stores.Count()
	meta["count"] = len(stores)
	meta["start"] = offset
	meta["sortBy"] = sb

	response := make(map[string]interface{})
	response["meta"] = meta
	response["items"] = stores

	app.jsonResponse(w, response)
}

func (app *application) getStore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	store, err := app.stores.Get(id)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, store)
}

// Flavor handlers
func (app *application) createFlavor(w http.ResponseWriter, r *http.Request) {
	var flavor = &models.Flavor{}
	err := json.NewDecoder(r.Body).Decode(&flavor)

	if err != nil {
		app.serverError(w, err)
	}
	defer r.Body.Close()

	flavor, err = app.flavors.Insert(flavor)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, flavor)
}

func (app *application) getFlavor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	flavor, err := app.flavors.Get(id)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, flavor)
}

func (app *application) listFlavor(w http.ResponseWriter, r *http.Request) {
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

	flavors, err := app.flavors.List(limit, offset, sb)

	if err != nil {
		app.serverError(w, err)
	}

	meta := make(map[string]interface{})
	meta["totalRecords"] = app.flavors.Count()
	meta["count"] = len(flavors)
	meta["start"] = offset
	meta["sortBy"] = sb

	response := make(map[string]interface{})
	response["meta"] = meta
	response["items"] = flavors

	app.jsonResponse(w, response)
}
