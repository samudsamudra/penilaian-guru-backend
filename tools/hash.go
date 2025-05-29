package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "jekingnigga"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	fmt.Println("Password:", password)
	fmt.Println("Hash:", string(hash))
}
