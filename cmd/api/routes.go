package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	// User routes
	mux.Post("/api/user", http.HandlerFunc(app.createUser))
	mux.Get("/api/user/:id", http.HandlerFunc(app.getUser))
	mux.Patch("/api/user/:id", http.HandlerFunc(app.partialUpdateUser))
	mux.Get("/api/user", http.HandlerFunc(app.listUser))
	mux.Del("/api/user/:id", http.HandlerFunc(app.deleteUser))

	// Store routes
	mux.Get("/api/store", http.HandlerFunc(app.listStore))
	mux.Post("/api/store", http.HandlerFunc(app.createStore))
	mux.Patch("/api/store/:id", http.HandlerFunc(app.partialUpdateStore))
	mux.Put("/api/store/:id", http.HandlerFunc(app.updateStore))
	mux.Get("/api/store/:id", http.HandlerFunc(app.getStore))
	// delete store

	// Flavor routes
	mux.Post("/api/flavor", http.HandlerFunc(app.createFlavor))
	mux.Get("/api/flavor", http.HandlerFunc(app.listFlavor))
	mux.Get("/api/flavor/:id", http.HandlerFunc(app.getFlavor))

	return app.logRequest(mux)
}
