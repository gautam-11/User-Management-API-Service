package Routes

import (
	"encoding/json"
	"net/http"
	"user-management-api-service/middlewares"
	"user-management-api-service/modules"
	"user-management-api-service/schemas"

	"github.com/go-chi/chi"
)

func UserRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middlewares.JwtAuthentication)
	//Handler for login
	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		credentials := schemas.LoginUser{}
		json.NewDecoder(r.Body).Decode(&credentials)
		err := middlewares.LoginValidate(&credentials)
		if err != nil {
			respondWithError(w, 400, err.Error())
			return
		}
		token, err := modules.LoginUser(&credentials)
		if err != nil {
			respondWithError(w, 400, err.Error())
		} else {
			respondWithJSON(w, 200, token)
		}

	})
	//Handler for handling registration of user
	router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		payload := schemas.User{}
		json.NewDecoder(r.Body).Decode(&payload)
		err := middlewares.Validate(&payload)
		if err != nil {
			respondWithError(w, 400, err.Error())
			return
		}
		user, err := modules.RegisterUser(&payload)
		if err != nil {
			respondWithError(w, 400, err.Error())
		} else {
			respondWithJSON(w, 200, user)
		}
	})

	//Handler for getting users by their mailid
	router.Get("/getuser/{email}", func(w http.ResponseWriter, r *http.Request) {

		userid := chi.URLParam(r, "email")
		user, err := modules.GetUserById(userid)
		if err != nil {
			respondWithError(w, 400, err.Error())
		} else {
			respondWithJSON(w, 200, user)
		}
	})
	//Handler for fetching all users
	router.Get("/getusers", func(w http.ResponseWriter, r *http.Request) {

		users, err := modules.GetUsers()
		if err != nil {
			respondWithError(w, 400, err.Error())
		} else {
			respondWithJSON(w, 200, users)
		}
	})
	//Handling for deleting user using mailid
	router.Delete("/{email}", func(w http.ResponseWriter, r *http.Request) {
		userid := chi.URLParam(r, "email")
		msg, err := modules.DeleteUser(userid)

		if err != nil {
			respondWithError(w, 400, err.Error())
		} else {
			respondWithJSON(w, 200, msg)
		}
	})
	router.Put("/{email}", func(w http.ResponseWriter, r *http.Request) {

		payload := schemas.User{}
		json.NewDecoder(r.Body).Decode(&payload)
		err := middlewares.Validate(&payload)
		if err != nil {
			respondWithError(w, 400, err.Error())
			return
		}
		userid := chi.URLParam(r, "email")
		msg, err := modules.UpdateUser(userid, &payload)
		if err != nil {
			respondWithError(w, 400, err.Error())
		} else {
			respondWithJSON(w, 200, msg)
		}

	})

	return router
}

// RespondwithJSON write json response format
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
