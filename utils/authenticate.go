package utils

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// CustomClaims - structure of claims of a jwt token
type CustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// CreateHash -  Function for storing password in encrypted form in DB
func CreateHash(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(hash), nil

}

// CheckHashedPassword - Compare user entered password with stored hash in DB
func CheckHashedPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true

}

// GenerateToken - Function to generate token
func GenerateToken(email string, role string, secret string) (string, error) {

	claims := CustomClaims{
		email,
		role,
		jwt.StandardClaims{
			ExpiresAt: (time.Now().Unix() + 100000),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), claims)

	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}
