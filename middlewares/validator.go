package middlewares

import (
	"fmt"
	"regexp"
	"user-management-api-service/schemas"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Validate payload
func Validate(user *schemas.User) error {
	fmt.Println("validation")
	err := validation.Errors{
		"first_name":  validation.Validate(user.FirstName, validation.Required, is.Alpha, validation.Length(3, 20)),
		"middle_name": validation.Validate(user.MiddleName, is.Alpha, validation.Length(3, 20)),
		"last_name":   validation.Validate(user.LastName, is.Alpha, validation.Length(3, 20)),
		"email":       validation.Validate(user.Email, validation.Required, is.Email),
		"phone":       validation.Validate(user.Phone, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{10,12}$"))),
		"password":    validation.Validate(user.Password, validation.Required, validation.Match(regexp.MustCompile("^[A-Za-z0-9_#$!@.]{10,30}$"))),
		"role":        validation.Validate(user.Role, validation.Required, validation.In("user", "admin")),
		"gender":      validation.Validate(user.Gender, validation.Required, validation.In("male", "female")),
	}.Filter()
	fmt.Println(err)
	return err

}
