package Routes

import (
	"net/http"

	"github.com/go-chi/chi"
)

func UserRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to this api service"))
	})
	/*
		router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("You are registered"))
		})
		router.Post("/login", modules.LoginUser)
		router.Get("/getusers", modules.GetUsers)
		router.Get("/getuser/{id}", modules.GetUser)
	*/
	return router
}
