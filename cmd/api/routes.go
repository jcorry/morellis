package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	// Auth route
	mux.Post("/api/v1/auth", http.HandlerFunc(app.createAuth))

	// User routes
	mux.Post("/api/v1/user", app.jwtVerification(NewPermissionsCheck(http.HandlerFunc(app.createUser), []string{"user:write", "self:write"})))
	mux.Get("/api/v1/user/:uuid", app.jwtVerification(NewPermissionsCheck(http.HandlerFunc(app.getUser), []string{"user:read", "self:read"})))
	mux.Patch("/api/v1/user/:id", app.jwtVerification(NewPermissionsCheck(http.HandlerFunc(app.partialUpdateUser), []string{"user:write", "self:write"})))
	mux.Get("/api/v1/user", app.jwtVerification(NewPermissionsCheck(http.HandlerFunc(app.listUser), []string{"user:read", "self:read"})))
	mux.Del("/api/v1/user/:uuid", app.jwtVerification(NewPermissionsCheck(http.HandlerFunc(app.deleteUser), []string{"user:write", "self:write"})))
	mux.Get("/api/v1/user/:uuid/ingredient", app.jwtVerification(NewPermissionsCheck(http.HandlerFunc(app.listUserIngredient), []string{"user:read", "self:read"})))
	mux.Post("/api/v1/user/:uuid/ingredient", app.jwtVerification(NewPermissionsCheck(http.HandlerFunc(app.createUserIngredientAssociation), []string{"user:write", "self:write"})))
	mux.Del("/api/v1/user/:uuid/ingredient/:userIngredientID", app.jwtVerification(NewPermissionsCheck(http.HandlerFunc(app.deleteUserIngredientAssociation), []string{"user:write", "self:write"})))

	// Store routes
	mux.Get("/api/v1/store", app.jwtVerification(http.HandlerFunc(app.listStore)))
	mux.Post("/api/v1/store", app.jwtVerification(http.HandlerFunc(app.createStore)))
	mux.Patch("/api/v1/store/:id", app.jwtVerification(http.HandlerFunc(app.partialUpdateStore)))
	mux.Put("/api/v1/store/:id", app.jwtVerification(http.HandlerFunc(app.updateStore)))
	mux.Get("/api/v1/store/:id", app.jwtVerification(http.HandlerFunc(app.getStore)))
	mux.Post("/api/v1/store/:storeID/flavor/:flavorID", app.jwtVerification(http.HandlerFunc(app.activateStoreFlavor)))
	mux.Del("/api/v1/store/:storeID/flavor/:flavorID", app.jwtVerification(http.HandlerFunc(app.deactivateStoreFlavor)))

	// Flavor routes
	mux.Post("/api/v1/flavor", app.jwtVerification(http.HandlerFunc(app.createFlavor)))
	mux.Get("/api/v1/flavor", app.jwtVerification(http.HandlerFunc(app.listFlavor)))
	mux.Get("/api/v1/flavor/:id", app.jwtVerification(http.HandlerFunc(app.getFlavor)))

	mux.Get("/api/v1/ingredient", app.jwtVerification(http.HandlerFunc(app.listIngredient)))

	return app.logRequest(mux)
}
