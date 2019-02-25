package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jcorry/morellis/pkg/models"
	"github.com/jcorry/morellis/pkg/models/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    interface {
		Insert(string, string, string, string, string) (*models.User, error)
		Update(*models.User) (*models.User, error)
		Get(int) (*models.User, error)
		List(int, int, string) ([]*models.User, error)
		Delete(int) (bool, error)
		Count() int
		Authenticate(string, string) (*models.User, error)
	}
	stores interface {
		Insert(string, string, string, string, string, string, string, string, float64, float64) (*models.Store, error)
		Update(int, string, string, string, string, string, string, string, string, float64, float64) (*models.Store, error)
		Get(int) (*models.Store, error)
	}
	mapsApiKey string
}

func main() {
	addr := flag.String("addr", ":4001", "HTTP network address")
	dsn := flag.String("dsn", "morellis:E4j+#2G^8Pa=^Nn9@(127.0.0.1:33061)/morellis?parseTime=true", "MySQL DSN")
	mapsApiKey := flag.String("api_key", "AIzaSyDzOe0YI-sQXJHM9DMr7YEU5zCwhPBXFts", "Google maps geocoding api key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:   errorLog,
		infoLog:    infoLog,
		users:      &mysql.UserModel{DB: db},
		stores:     &mysql.StoreModel{DB: db},
		mapsApiKey: *mapsApiKey,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)

	err = srv.ListenAndServe()
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
