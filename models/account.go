package models

import (
	"os"
	"strings"

	"github.com/pkg/errors"

	u "github.com/jcorry/morellis/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	Base
	Email           string        `json:"email"`
	Password        string        `json:"password"`
	Token           string        `json:"token";sql:"-"`
	AccountStatusID uint          `json:"-"`
	AccountStatus   AccountStatus `json:"accountStatus"`
}

type AccountStatus struct {
	ID    uint   `json:"-"`
	Value string `json:"value";sql:"not null;type:ENUM('Pending','Active','Deleted')"`
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).Where("account_status_id = ?", 1).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() (*Account, error) {

	if _, ok := account.Validate(); !ok {
		return nil, errors.New("Invalid account")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return nil, errors.New("Unable to create new account")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	return account, nil
}

func Login(email, password string) (*Account, error) {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ? ", email).First(account).Error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, err
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	return account, nil
}

func GetUser(u uint) (*Account, error) {
	account := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(account)
	if account.Email == "" { //User not found!
		return nil, errors.New("User not found")
	}

	account.Password = ""
	return account, nil
}
