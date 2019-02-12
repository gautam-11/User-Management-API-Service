package tests

import (
	"fmt"
	"log"
	"testing"
	"user-management-api-service/utils"

	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJwt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JWT check test")
}

var _ = Describe("Token generation and validation", func() {

	Context("When user sends a jwt token", func() {
		It("should be a valid jwt token", func() {

			email := "gautam.b@gmail.com"
			role := "admin"
			configuration, conferr := GetEnv()
			if conferr != nil {
				log.Println("Configuration error", conferr)
			}

			JWTsecret := configuration.Constants.JWT_SECRET

			tokenPart, err := utils.GenerateToken(email, role, JWTsecret)

			if err != nil {
				log.Println("Error in generating token.....", err)
				return
			}

			tk := &utils.CustomClaims{}
			_, err = jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(JWTsecret), nil
			})
			if err != nil {
				log.Println("Error in decoding token ..", err)
				return
			}

			Expect(tk.Email).Should(Equal(email))

			Expect(tk.Role).Should(Equal(role))

			fmt.Println(tk.Email, tk.Role)
		})
	})

})
