package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jcorry/morellis/pkg/models"
	"github.com/jcorry/morellis/pkg/models/mysql"
	"github.com/jcorry/morellis/pkg/sms"
)

type UserIngredientBody struct {
	ID           int64     `json:"id"`
	UserUUID     uuid.UUID `json:"userUuid"`
	IngredientID int64     `json:"ingredientId"`
	StoreID      int64     `json:"storeId,omitempty"`
	Keyword      string    `json:"keyword,omitempty"`
	Created      time.Time `json:"created"`
}

// Webhook handlers

// smsAuthRequest looks the user up by their phone number, supplied by the incoming twilio
// webhook. If found, generates an expiring auth token and sends the user a URL at
// which they can authenticate and get a JWT with limited permissions for future requests
func (app *application) smsAuthRequest(w http.ResponseWriter, r *http.Request) {
	err := sms.ValidateIncomingRequest(app.baseUrl, os.Getenv("TWILIO_AUTH_TOKEN"), r)
	if err != nil {
		app.errorLog.Output(2, err.Error())
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	guid := uuid.New()
	token := base64.StdEncoding.EncodeToString([]byte(guid.String()))

	type TwilioPayload struct {
		From string `json:"From"`
		Body string `json:"Body"`
	}
	var tp TwilioPayload
	err = json.NewDecoder(r.Body).Decode(&tp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user *models.User
	user, err = app.users.GetByPhone(tp.From)
	// If no user is found, create one
	if err != nil {
		if err == models.ErrNoRecord {
			password, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			user, err = app.users.Insert(uuid.New(), models.NullString{}, models.NullString{}, models.NullString{}, tp.From, int(models.USER_STATUS_VERIFIED), string(password))
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// save the token
	err = app.users.SaveAuthToken(token, int(user.ID))
	if err != nil {
		app.serverError(w, err)
		return
	}

	url := fmt.Sprintf(`https://%s/auth/%s`, app.baseUrl, token)
	message := fmt.Sprintf(`access the 🍦 app at: %s`, url)

	_, err = app.sender.Send(r.Context(), user.Phone, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// authByToken looks for a valid auth token in the URL and if found, returns a JWT
func (app *application) authByToken(w http.ResponseWriter, r *http.Request) {
	// look up user by token
	user, err := app.users.GetByAuthToken(r.URL.Query().Get(":token"))
	if err != nil {
		app.errorLog.Output(2, err.Error())
		if err == models.ErrNoRecord {
			app.clientError(w, http.StatusNotFound)
			return
		}
		if err == mysql.ErrNoAuthTokenFound {
			app.clientError(w, http.StatusNotFound)
			return
		}
		app.serverError(w, err)
		return
	}
	user.Permissions, err = app.users.GetPermissions(int(user.ID))
	if err != nil {
		app.serverError(w, err)
		return
	}

	token, err := generateToken(user)
	if err != nil {
		app.serverError(w, err)
		return
	}

	claims, err := verifyToken(token)
	if err != nil {
		app.serverError(w, err)
		return
	}

	exp := time.Unix(claims.ExpiresAt, 0)

	response := struct {
		Token   string    `json:"token"`
		Expires time.Time `json:"expires"`
	}{
		token,
		exp,
	}

	app.jsonResponse(w, response)

}

// Auth handlers
func (app *application) createAuth(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer r.Body.Close()

	user, err := app.users.GetByCredentials(creds)
	if err != nil {
		app.errorLog.Output(2, err.Error())
		app.clientError(w, http.StatusNotFound)
		return
	}

	user.Permissions, err = app.users.GetPermissions(int(user.ID))
	if err != nil {
		app.serverError(w, err)
		return
	}

	token, err := generateToken(user)
	if err != nil {
		app.serverError(w, err)
		return
	}

	claims, err := verifyToken(token)
	if err != nil {
		app.serverError(w, err)
		return
	}

	exp := time.Unix(claims.ExpiresAt, 0)

	response := struct {
		Token   string    `json:"token"`
		Expires time.Time `json:"expires"`
	}{
		token,
		exp,
	}

	app.jsonResponse(w, response)
}

// User handlers
func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var reqUser *models.User
	err := json.NewDecoder(r.Body).Decode(&reqUser)

	if err != nil {
		app.serverError(w, err)
		return
	}
	defer r.Body.Close()

	uid, err := uuid.NewRandom()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var userStatus models.UserStatus

	var user *models.User
	user, err = app.users.Insert(uid, reqUser.FirstName, reqUser.LastName, reqUser.Email, reqUser.Phone, int(userStatus.GetID(reqUser.Status)), reqUser.Password)

	if err != nil {
		if err == models.ErrDuplicateEmail || err == models.ErrDuplicatePhone {
			app.badRequest(w, err)
			return
		}

		app.serverError(w, err)
		return
	}
	user.Password = ""
	user.UUID = uid

	for _, up := range reqUser.Permissions {
		_, err = app.users.AddPermission(int(user.ID), up.Permission)
		if err != nil {
			app.serverError(w, err)
			return
		}
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
		return
	}

	user, err = app.users.Update(user)
	if err != nil {
		if err == models.ErrDuplicateEmail || err == models.ErrDuplicatePhone {
			app.badRequest(w, err)
			return
		}

		app.serverError(w, err)
		return
	}

	app.jsonResponse(w, user)
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get(":uuid"))
	if err != nil || id == uuid.Nil {
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

	user.Permissions, err = app.users.GetPermissions(int(user.ID))
	if err != nil {
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
		return
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
	id, err := uuid.Parse(r.URL.Query().Get(":uuid"))
	if err != nil || id == uuid.Nil {
		app.notFound(w)
		return
	}

	user, err := app.users.GetByUUID(id)
	if err != nil {
		app.notFound(w)
		return
	}

	userID := int(user.ID)
	// Remove all of the User Permissions
	err = app.users.RemoveAllPermissions(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	res, err := app.users.Delete(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if res {
		app.noContentResponse(w)
		return
	}
}

func (app *application) createUserIngredientAssociation(w http.ResponseWriter, r *http.Request) {
	userUUID, err := uuid.Parse(r.URL.Query().Get(":uuid"))
	if err != nil || userUUID == uuid.Nil {
		fmt.Println("No user found")
		app.notFound(w)
		return
	}

	user, err := app.users.GetByUUID(userUUID)
	if err != nil {
		app.notFound(w)
		return
	}

	var userIngredient UserIngredientBody

	err = json.NewDecoder(r.Body).Decode(&userIngredient)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If that's not an actual Ingredient: 404
	ingredient, err := app.ingredients.Get(userIngredient.IngredientID)
	if err != nil {
		app.infoLog.Output(2, fmt.Sprintf("No ingredient found for ID: %d", userIngredient.IngredientID))
		app.notFound(w)
		return
	}

	ui, err := app.users.AddIngredient(user.ID, ingredient, userIngredient.Keyword)
	if err != nil {
		if err == models.ErrDuplicateUserIngredient {
			app.clientError(w, http.StatusConflict)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}
	userIngredient.ID = ui.UserIngredientID
	userIngredient.Created = ui.Created

	app.jsonResponse(w, &userIngredient)
}

func (app *application) deleteUserIngredientAssociation(w http.ResponseWriter, r *http.Request) {
	userUUID, err := uuid.Parse(r.URL.Query().Get(":uuid"))
	if err != nil || userUUID == uuid.Nil {
		fmt.Println("No user found")
		app.notFound(w)
		return
	}

	_, err = app.users.GetByUUID(userUUID)
	if err != nil {
		app.notFound(w)
		return
	}

	userIngredientId, err := strconv.Atoi(r.URL.Query().Get(":userIngredientID"))
	if err != nil {
		app.infoLog.Output(2, "No userIngredientID found in URL")
		app.notFound(w)
		return
	}

	err = app.users.RemoveUserIngredient(int64(userIngredientId))

	if err != nil {
		app.errorLog.Output(2, err.Error())

		if err == models.ErrNoneAffected {
			app.notFound(w)
			return
		}

		app.serverError(w, err)
		return
	}

	app.noContentResponse(w)
	return
}

func (app *application) listUserIngredient(w http.ResponseWriter, r *http.Request) {
	userUUID, err := uuid.Parse(r.URL.Query().Get(":uuid"))
	if err != nil || userUUID == uuid.Nil {
		fmt.Println("No user found")
		app.notFound(w)
		return
	}

	user, err := app.users.GetByUUID(userUUID)
	if err != nil {
		app.notFound(w)
		return
	}
	fmt.Println(fmt.Sprintf("User UUID: %s", userUUID))
	fmt.Println(fmt.Sprintf("UserID: %d", user.ID))

	userIngredients, err := app.users.GetIngredients(user.ID)
	if err != nil {
		app.notFound(w)
		return
	}

	userIngredientResponses := []*UserIngredientBody{}

	for _, ui := range userIngredients {
		userIngredientResponses = append(userIngredientResponses, &UserIngredientBody{
			ID:           ui.UserIngredientID,
			UserUUID:     userUUID,
			IngredientID: ui.Ingredient.ID,
			Created:      ui.Created,
		})
	}

	meta := make(map[string]interface{})
	meta["totalRecords"] = len(userIngredientResponses)
	meta["count"] = len(userIngredientResponses)

	response := make(map[string]interface{})
	response["meta"] = meta
	response["items"] = userIngredientResponses

	app.jsonResponse(w, response)
	return
}

// Store handlers
func (app *application) createStore(w http.ResponseWriter, r *http.Request) {
	var store *models.Store
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = json.Unmarshal(b, &store)

	if err != nil {
		app.serverError(w, err)
		return
	}

	// Geocode the store
	err = app.geocodeStore(store)
	if err != nil {
		app.serverError(w, err)
		return
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
		return
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

func (app *application) activateStoreFlavor(w http.ResponseWriter, r *http.Request) {
	storeID, err := strconv.Atoi(r.URL.Query().Get(":storeID"))
	if err != nil || storeID < 1 {
		app.notFound(w)
		return
	}

	flavorID, err := strconv.Atoi(r.URL.Query().Get(":flavorID"))
	if err != nil || flavorID < 1 {
		app.notFound(w)
		return
	}

	type activationRequestBody struct {
		StoreID  int64     `json:"store_id"`
		FlavorID int64     `json:"flavor_id"`
		Position int       `json:"position"`
		Created  time.Time `json:"created,omitempty"`
	}

	var req activationRequestBody

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		app.serverError(w, err)
		return
	}

	if req.FlavorID != int64(flavorID) {
		app.errorLog.Output(2, fmt.Sprintf("Request body flavor_id (%d) must match URL query :flavorID (%d)", req.FlavorID, flavorID))
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if req.StoreID != int64(storeID) {
		app.errorLog.Output(2, fmt.Sprintf("Request body store_id (%d) must match URL query :storeID (%d)", req.StoreID, storeID))
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Is it even an available Flavor?
	f, err := app.flavors.Get(flavorID)
	if err == models.ErrNoRecord {
		app.clientError(w, http.StatusNotFound)
		return
	}
	req.FlavorID = f.ID

	// Is it a valid store?
	s, err := app.stores.Get(storeID)
	if err == models.ErrNoRecord {
		app.clientError(w, http.StatusNotFound)
		return
	}
	req.StoreID = s.ID

	// Make the association link
	err = app.stores.ActivateFlavor(s.ID, f.ID, req.Position)
	if err != nil {
		app.serverError(w, err)
		return
	}
	req.Created = time.Now()

	app.jsonResponse(w, req)
}

func (app *application) deactivateStoreFlavor(w http.ResponseWriter, r *http.Request) {
	storeID, err := strconv.Atoi(r.URL.Query().Get(":storeID"))
	if err != nil || storeID < 1 {
		app.notFound(w)
		return
	}

	flavorID, err := strconv.Atoi(r.URL.Query().Get(":flavorID"))
	if err != nil || flavorID < 1 {
		app.notFound(w)
		return
	}

	_, err = app.stores.DeactivateFlavor(int64(storeID), int64(flavorID))
	if err != nil {
		app.errorLog.Output(2, err.Error())
		app.clientError(w, http.StatusBadRequest)
	}

	app.noContentResponse(w)
}

// Flavor handlers
func (app *application) createFlavor(w http.ResponseWriter, r *http.Request) {
	var flavor = &models.Flavor{}
	err := json.NewDecoder(r.Body).Decode(&flavor)

	if err != nil {
		app.serverError(w, err)
		return
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
	fi := params.Get("filterIngredient")

	t := csv.NewReader(strings.NewReader(fi))

	var ingredientTerms []string
	for {
		r, err := t.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			app.serverError(w, err)
			return
		}
		ingredientTerms = r
	}

	flavors, err := app.flavors.List(limit, offset, sb, ingredientTerms)

	if err != nil {
		app.serverError(w, err)
		return
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

// Ingredient handlers
func (app *application) listIngredient(w http.ResponseWriter, r *http.Request) {
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

	s := params.Get("searchTerms")

	t := csv.NewReader(strings.NewReader(s))

	var terms []string
	for {
		r, err := t.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			app.serverError(w, err)
			return
		}
		terms = r
	}

	ingredients, err := app.ingredients.Search(limit, offset, sb, terms)
	if err != nil {
		app.serverError(w, err)
		return
	}

	meta := make(map[string]interface{})
	meta["totalRecords"] = app.flavors.Count()
	meta["count"] = len(ingredients)
	meta["start"] = offset
	meta["sortBy"] = sb

	response := make(map[string]interface{})
	response["meta"] = meta
	response["items"] = ingredients

	app.jsonResponse(w, response)
}
