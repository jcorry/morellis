package models

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var db *gorm.DB

type Base struct {
	ID        uint       `json:"id,omitempty";gorm:"primary_key"`
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated,omitempty"`
	DeletedAt *time.Time `json:"deleted,omitempty";sql:"index"`
}

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbPort := os.Getenv("db_port")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbHost, dbPort, dbName)
	fmt.Println(dbUri)

	conn, err := gorm.Open("mysql", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Set("association_autoupdate", true)
	db.Debug().AutoMigrate(&Account{}, &AccountStatus{}, &Flavor{}, &Ingredient{})
}

func GetDB() *gorm.DB {
	return db
}
