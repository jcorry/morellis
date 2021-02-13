package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/cors"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jcorry/morellis/pkg/models"
	"github.com/jcorry/morellis/pkg/models/mysql"
	"github.com/jcorry/morellis/pkg/sms"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	users       models.UserRepository
	stores      models.StoreRepository
	flavors     models.FlavorRepository
	ingredients models.IngredientRepository
	sender      sms.Messager
	baseUrl     string
	mapsApiKey  string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	mapsApiKey := os.Getenv("GMAP_API_KEY")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)

	if err != nil {
		errorLog.Fatal(err)
	}
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(101)

	defer db.Close()

	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Initialize Twilio Client
	client := &http.Client{}
	sender := sms.NewTwilioMessager(client, os.Getenv("TWILIO_SID"), os.Getenv("TWILIO_AUTH_TOKEN"), os.Getenv("TWILIO_NUMBER"))

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		users:       &mysql.UserModel{DB: db, Redis: rdb},
		stores:      &mysql.StoreModel{DB: db},
		flavors:     &mysql.FlavorModel{DB: db},
		ingredients: &mysql.IngredientModel{DB: db},
		mapsApiKey:  mapsApiKey,
		sender:      sender,
		baseUrl:     os.Getenv("HOST"),
	}

	c := cors.New(cors.Options{
		AllowedOrigins:     []string{fmt.Sprintf("%s:*", os.Getenv("HOST"))},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   true,
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPatch, http.MethodDelete, http.MethodPut},
		OptionsPassthrough: true,
		Debug:              true,
	})

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errorLog,
		Handler:      c.Handler(app.routes()),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	infoLog.Printf("Starting server on %s", addr)

	err = srv.ListenAndServeTLS(`./tls/cert.pem`, `./tls/key.pem`)
	errorLog.Fatal(err)
}

// openDB opens a DB connection using for a dsn
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
