package utils

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

//Generate Hash from password to store in DB
func CreateHash(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(hash), nil

}

//Compare user entered password with stored hash in DB
func CheckHashedPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true

}

func GenerateToken(email string, role string, secret string) (string, error) {

	claims := CustomClaims{
		email,
		role,
		jwt.StandardClaims{
			ExpiresAt: (time.Now().Unix() + 1000),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), claims)

	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}
