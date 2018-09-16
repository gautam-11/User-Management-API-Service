package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"user-management-api-service/internal/config"
	"user-management-api-service/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestPath := r.URL.Path //current request path
		allowedRoutes := make(map[string]bool)
		allowedRoutes["/v1/api/user/login"] = true
		allowedRoutes["/v1/api/user/register"] = true

		//check if request does not need authentication, serve the request if it doesn't need it
		if allowedRoutes[requestPath] == true {
			next.ServeHTTP(w, r)
			return
		}

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			respondWithError(w, http.StatusForbidden, "Missing auth token")
			return
		}

		splitted := strings.Split(tokenHeader, " ") //token format `Bearer {token-body}`
		if len(splitted) != 2 {
			respondWithError(w, http.StatusForbidden, "Invalid/Malformed auth token")
			return
		}

		tokenPart := splitted[1] //Grab the token part
		tk := &utils.CustomClaims{}

		secrets, err := config.GetEnv()
		if err != nil {
			log.Panicln("Configuration error", err)
			return
		}
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(secrets.Constants.JWT_SECRET), nil
		})

		fmt.Println(err)
		if err != nil { //Malformed token, returns with http code 403 as usual
			respondWithError(w, http.StatusForbidden, "Malformed Token")
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			respondWithError(w, http.StatusForbidden, "Invalid token")
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Println("Email ", tk.Email)
		//ctx := context.WithValue(r.Context(), "user", tk.UserId)
		//r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"message": msg})
}