package schemas

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName  string        `json:"first_name" bson:"first_name"`
	MiddleName string        `json:"middle_name,omitempty" bson:"middle_name"`
	LastName   string        `json:"last_name" bson:"last_name"`
	Email      string        `json:"email" bson:"email"`
	Phone      string        `json:"phone,omitempty" bson:"phone"`
	Password   string        `json:"password" bson:"password"`
	Role       string        `json:"role,omitempty" bson:"role"`
	Gender     string        `json:"gender,omitempty" bson:"gender"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
}
