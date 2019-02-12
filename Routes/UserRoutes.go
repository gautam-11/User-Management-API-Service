package Routes

import (
	"encoding/json"
	"net/http"
	"user-management-api-service/middlewares"
	"user-management-api-service/modules"
	"user-management-api-service/schemas"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// SetRoutes is a method for setting up the router - returns a handler
func SetRoutes() *chi.Mux {
	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,                             // Log API request calls
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
		cors.Handler,
	)

	//Root path testing
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to User management api service"))
	})

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/user", UserRoutes())
	})
	return router
}

// UserRoutes is a method for defining all the  CRUD api endpoints
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
			return
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
			return
		} else {
			respondWithJSON(w, 200, user)
		}
	})

	//Handler for getting users by their mailid
	router.Get("/getuser/{email}", func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)

		CanAccess := RoleChecker(role, 1) // 1 -> Read operation
		if !CanAccess {
			respondWithError(w, 400, "You are not authorized to perform this operation")
			return
		}
		userid := chi.URLParam(r, "email")
		user, err := modules.GetUserById(userid)
		if err != nil {
			respondWithError(w, 400, err.Error())
			return
		} else {
			respondWithJSON(w, 200, user)
		}
	})
	//Handler for fetching all users
	router.Get("/getusers", func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value("role").(string)
		CanAccess := RoleChecker(role, 1) // 1 -> Read operation
		if !CanAccess {
			respondWithError(w, 400, "You are not authorized to perform this operation")
			return
		}
		users, err := modules.GetUsers()
		if err != nil {
			respondWithError(w, 400, err.Error())
		} else {
			respondWithJSON(w, 200, users)
		}
	})
	//Handling for deleting user using mailid
	router.Delete("/{email}", func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		CanAccess := RoleChecker(role, 3) // 3 -> Delete operation
		if !CanAccess {
			respondWithError(w, 400, "You are not authorized to perform this operation")
			return
		}
		userid := chi.URLParam(r, "email")
		msg, err := modules.DeleteUser(userid)

		if err != nil {
			respondWithError(w, 400, err.Error())
		} else {
			respondWithJSON(w, 200, msg)
		}
	})
	router.Put("/{email}", func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value("role").(string)
		CanAccess := RoleChecker(role, 2) // 2 -> Update operation
		if !CanAccess {
			respondWithError(w, 400, "You are not authorized to perform this operation")
			return
		}

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
			return
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

// RoleChecker : Maps a user's role to its allowed operation
func RoleChecker(role string, crud uint8) bool {

	return modules.RoleMap(role, crud)

}
