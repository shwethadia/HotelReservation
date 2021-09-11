package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHash() {

	password := "1234"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	fmt.Println(string(hashedPassword))
}


