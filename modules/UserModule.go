package modules

import (
	"errors"
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

	resp := utils.UserJson{payload.Email, payload.Phone, "User Registered successfully"}

	return resp, err
}
