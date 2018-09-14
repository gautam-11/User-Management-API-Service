package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CreateHash(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(hash), nil

}
