package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "password"

	err := bcrypt.CompareHashAndPassword([]byte("$2a$14$f3CbQRltceEckV2GkSTiCu3XjyJ088Q9oL52gWMWMUGjmyt3lwxCG"), []byte(password))

	if err != nil {
		fmt.Println(fmt.Errorf("Password and hash don't match"))
	} else {
		fmt.Print("Passwords match!")
	}
}
