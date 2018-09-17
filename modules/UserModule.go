package modules

import (
	"errors"
	"fmt"
	"log"
	"time"
	"user-management-api-service/internal/config"
	"user-management-api-service/schemas"
	"user-management-api-service/utils"

	"gopkg.in/mgo.v2/bson"
)

// RegisterUser DB connection logic
func RegisterUser(payload *schemas.User) (utils.UserJson, error) {
	db, err := config.Connect()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	defer db.Session.Close()

	var result schemas.User
	err = db.Database.C("users").Find(bson.M{"$or": []bson.M{bson.M{"email": payload.Email}, bson.M{"phone": payload.Phone}}}).One(&result)

	if err == nil {
		return utils.UserJson{}, errors.New("Email or Phone already exists")
	}
	payload.CreatedAt = time.Now()
	payload.Password, _ = utils.CreateHash(payload.Password)
	err = db.Database.C("users").Insert(payload)
	if err != nil {
		return utils.UserJson{}, err
	}

	resp := utils.UserJson{payload.FirstName, payload.LastName, payload.Email, payload.Phone, "User Registered successfully"}

	return resp, err
}

//Login User DB connection login

func LoginUser(credentials *schemas.LoginUser) (string, error) {
	db, err := config.Connect()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	defer db.Session.Close()

	var result *schemas.User

	err = db.Database.C("users").Find(bson.M{"email": credentials.Email}).One(&result)

	if err != nil {
		return "", errors.New("User not found!!")
	}
	if !(utils.CheckHashedPassword(result.Password, credentials.Password)) {
		return "", errors.New("Invalid Password")
	}
	token, err := utils.GenerateToken(result.Email, result.Role, db.Constants.JWT_SECRET)
	if err != nil {
		return "", err
	}
	return token, nil

}

//Search for User using mailid / phone

func GetUserById(email string) (utils.UserJson, error) {
	db, err := config.Connect()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	defer db.Session.Close()
	var result schemas.User
	err = db.Database.C("users").Find(bson.M{"email": email}).One(&result)

	if err != nil {
		return utils.UserJson{}, errors.New("User not found!!")
	}
	resp := utils.UserJson{result.FirstName, result.LastName, result.Email, result.Phone, "Fetched user successfully"}
	return resp, err

}

//Fetch all users

func GetUsers() ([]schemas.User, error) {
	db, err := config.Connect()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	defer db.Session.Close()
	var results []schemas.User
	err = db.Database.C("users").Find(bson.M{}).Select(bson.M{"password": 0, "_id": 0}).All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil

}

// Delete user using mailid
func DeleteUser(email string) (string, error) {
	db, err := config.Connect()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	defer db.Session.Close()
	err = db.Database.C("users").Remove(bson.M{"email": email})
	if err != nil {
		return "", err
	}

	return "Deleted user successfully", nil
}

//Update user using mailid

func UpdateUser(email string, payload *schemas.User) (string, error) {
	db, err := config.Connect()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	defer db.Session.Close()

	payload.UpdatedAt = time.Now()
	err = db.Database.C("users").Update(bson.M{"email": email}, payload)
	if err != nil {
		return "", err
	}

	return "Updated user successfully", nil

}

//Check for Existence of User based on token claims
func DoesUserExist(email string, role string) bool {
	db, err := config.Connect()
	if err != nil {

		log.Panicln("Configuration error", err)
	}
	defer db.Session.Close()

	var result schemas.User
	err = db.Database.C("users").Find(bson.M{"$and": []bson.M{bson.M{"email": email}, bson.M{"role": role}}}).One(&result)

	fmt.Println("Error", err)
	if err == nil {
		return true
	}

	return false
}
