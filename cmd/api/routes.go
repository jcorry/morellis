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
	mux.Post("/api/v1/user", app.jwtVerification(http.HandlerFunc(app.createUser)))
	mux.Get("/api/v1/user/:uuid", app.jwtVerification(http.HandlerFunc(app.getUser)))
	mux.Patch("/api/v1/user/:id", app.jwtVerification(http.HandlerFunc(app.partialUpdateUser)))
	mux.Get("/api/v1/user", app.jwtVerification(http.HandlerFunc(app.listUser)))
	mux.Del("/api/v1/user/:id", app.jwtVerification(http.HandlerFunc(app.deleteUser)))

	// Store routes
	// list stores
	mux.Get("/api/v1/store", app.jwtVerification(http.HandlerFunc(app.listStore)))
	mux.Post("/api/v1/store", app.jwtVerification(http.HandlerFunc(app.createStore)))
	mux.Patch("/api/v1/store/:id", app.jwtVerification(http.HandlerFunc(app.partialUpdateStore)))
	mux.Put("/api/v1/store/:id", app.jwtVerification(http.HandlerFunc(app.updateStore)))
	mux.Get("/api/v1/store/:id", app.jwtVerification(http.HandlerFunc(app.getStore)))
	mux.Post("/api/v1/store/:storeID/flavor/:flavorID", app.jwtVerification(http.HandlerFunc(app.activateStoreFlavor)))
	mux.Del("/api/v1/store/:storeID/flavor/:flavorID", app.jwtVerification(http.HandlerFunc(app.deactivateStoreFlavor)))
	// delete store

	// Flavor routes
	mux.Post("/api/v1/flavor", app.jwtVerification(http.HandlerFunc(app.createFlavor)))
	mux.Get("/api/v1/flavor", app.jwtVerification(http.HandlerFunc(app.listFlavor)))
	mux.Get("/api/v1/flavor/:id", app.jwtVerification(http.HandlerFunc(app.getFlavor)))

	return app.logRequest(mux)
}
