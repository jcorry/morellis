// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jcorry/morellis/pkg/nlp"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
	}
	var dsn = fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	db, err := openDB(dsn)

	svc, err := nlp.NewEntityTypesService()
	if err != nil {
		fmt.Errorf("%s", err)
		log.Fatal(err)
	}
	svc.DB = db

	// Add the flavors and ingredients
	err = svc.AddEntities()
	if err != nil {
		fmt.Errorf("%s", err)
	}

	// Get new list
	types, err := svc.ListEntityTypes()
	if err != nil {
		fmt.Errorf("%v", err)
	}

	// List them again
	fmt.Println(fmt.Sprintf("Line 28: %+v", types))
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
