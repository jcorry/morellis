package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/getsentry/raven-go"

	"github.com/jcorry/morellis/app"
	"github.com/jcorry/morellis/controllers"

	"github.com/gorilla/mux"
)

func main() {

	raven.SetDSN("https://79c1a519cb1f4bbca522c36c2fcaf975:f984b3bc27844e6196f54452451f74a1@sentry.io/1322019")

	_, err := os.Open("filename.ext")
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/flavor", controllers.CreateFlavor).Methods("POST")
	router.HandleFunc("/api/flavor", controllers.GetFlavors).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("app_port")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err = http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
