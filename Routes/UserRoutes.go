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

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to User management api service"))
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
	/*
		router.Post("/login", modules.LoginUser)
		router.Get("/getusers", modules.GetUsers)
		router.Get("/getuser/{id}", modules.GetUser)
	*/
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
