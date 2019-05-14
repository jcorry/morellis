package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jcorry/morellis/pkg/models"
	"github.com/jcorry/morellis/pkg/models/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    interface {
		Insert(uid uuid.UUID, firstName string, lastName string, email string, phone string, statusID int, password string) (*models.User, error)
		Update(*models.User) (*models.User, error)
		Get(int) (*models.User, error)
		GetByUUID(uuid.UUID) (*models.User, error)
		GetByCredentials(models.Credentials) (*models.User, error)
		List(int, int, string) ([]*models.User, error)
		Delete(int) (bool, error)
		Count() int
		GetPermissions(userID int) ([]models.UserPermission, error)
		AddPermission(userID int, p models.Permission) (int, error)
		RemovePermission(userPermissionID int) (bool, error)
		RemoveAllPermissions(userID int) error
	}
	stores interface {
		Insert(string, string, string, string, string, string, string, string, float64, float64) (*models.Store, error)
		Update(int, string, string, string, string, string, string, string, string, float64, float64) (*models.Store, error)
		Get(storeID int) (*models.Store, error)
		List(int, int, string) ([]*models.Store, error)
		Count() int
		ActivateFlavor(storeID int64, flavorID int64, position int) error
		DeactivateFlavor(storeID int64, flavorID int64) (bool, error)
		DeactivateFlavorAtPosition(storeID int64, position int) (bool, error)
	}
	flavors interface {
		Count() int
		Get(int) (*models.Flavor, error)
		List(int, int, string) ([]*models.Flavor, error)
		Insert(*models.Flavor) (*models.Flavor, error)
		Update(int, *models.Flavor) (*models.Flavor, error)
		Delete(int) (bool, error)
		AddIngredient(int, *models.Ingredient) (*models.Ingredient, error)
		RemoveIngredient(int, *models.Ingredient) (*models.Ingredient, error)
	}
	ingredients interface {
		GetByName(string) (*models.Ingredient, error)
		Insert(*models.Ingredient) (*models.Ingredient, error)
		Search(limit int, offset int, order string, search []string) ([]*models.Ingredient, error)
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
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(101)

	defer db.Close()

	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		users:       &mysql.UserModel{DB: db},
		stores:      &mysql.StoreModel{DB: db},
		flavors:     &mysql.FlavorModel{DB: db},
		ingredients: &mysql.IngredientModel{DB: db},
		mapsApiKey:  *mapsApiKey,
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
