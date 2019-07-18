package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"

	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jcorry/morellis/pkg/models"
	"github.com/jcorry/morellis/pkg/models/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    interface {
		Insert(uid uuid.UUID, firstName models.NullString, lastName models.NullString, email models.NullString, phone string, statusID int, password string) (*models.User, error)
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
		AddIngredient(userID int64, ingredient *models.Ingredient, keyword string) (*models.UserIngredient, error)
		GetIngredients(userID int64) ([]*models.UserIngredient, error)
		RemoveUserIngredient(userIngredientID int64) error
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
		List(limit int, offset int, sortBy string, ingredientTerms []string) ([]*models.Flavor, error)
		Insert(*models.Flavor) (*models.Flavor, error)
		Update(int, *models.Flavor) (*models.Flavor, error)
		Delete(int) (bool, error)
	}
	ingredients interface {
		Get(ID int64) (*models.Ingredient, error)
		GetByName(string) (*models.Ingredient, error)
		Insert(*models.Ingredient) (*models.Ingredient, error)
		Search(limit int, offset int, order string, search []string) ([]*models.Ingredient, error)
	}
	mapsApiKey string
}

type MySQLConfig struct {
	// Optional.
	Username, Password string

	// Database name, used for separating local from GKE envs
	Database string

	// Host of the MySQL instance.
	//
	// If set, UnixSocket should be unset.
	Host string

	// Port of the MySQL instance.
	//
	// If set, UnixSocket should be unset.
	Port int

	// UnixSocket is the filepath to a unix socket.
	//
	// If set, Host and Port should be unset.
	UnixSocket string
}

func main() {
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	// dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))

	mapsApiKey := os.Getenv("GMAP_API_KEY")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Println("Connecting to DB...")

	db, err := configureCloudSQL(cloudSQLConfig{
		Username: os.Getenv("GCP_MYSQL_USERNAME"),
		Password: os.Getenv("GCP_MYSQL_PASSWORD"),
		Instance: os.Getenv("GCP_MYSQL_INSTANCE"),
	})

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
		mapsApiKey:  mapsApiKey,
	}

	infoLog.Println("Setting up CORS headers...")
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:*"},
		AllowedHeaders: []string{"Authorization", "X-Requested-With", "Content-Type"},
		Debug:          true,
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

type cloudSQLConfig struct {
	Username, Password, Instance string
}

func configureCloudSQL(config cloudSQLConfig) (*sql.DB, error) {
	if os.Getenv("GAE_INSTANCE") != "" {
		// Running in production.
		return newMySQLDB(MySQLConfig{
			Username:   config.Username,
			Password:   config.Password,
			UnixSocket: "/cloudsql/" + config.Instance,
		})
	}

	// Running locally.
	return newMySQLDB(MySQLConfig{
		Username: os.Getenv("GCP_MYSQL_USERNAME"),
		Password: os.Getenv("GCP_MYSQL_PASSWORD"),
		Host:     "127.0.0.1",
		Database: os.Getenv("GCP_MYSQL_DATABASE"),
		Port:     3306,
	})
}

// newMySQLDB creates a new database backed by a given MySQL server.
func newMySQLDB(config MySQLConfig) (*sql.DB, error) {
	conn, err := sql.Open("mysql", config.dataStoreName(config.Database))
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
	}
	return conn, nil
}

// dataStoreName returns a connection string suitable for sql.Open.
func (c MySQLConfig) dataStoreName(databaseName string) string {
	var cred string
	// [username[:password]@]
	if c.Username != "" {
		cred = c.Username
		if c.Password != "" {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	if c.UnixSocket != "" {
		return fmt.Sprintf("%sunix(%s)/%s?parseTime=true", cred, c.UnixSocket, databaseName)
	}
	return fmt.Sprintf("%stcp([%s]:%d)/%s?parseTime=true", cred, c.Host, c.Port, databaseName)
}
