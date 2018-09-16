package main

import (
	"log"
	"net/http"

	"user-management-api-service/Routes"
	"user-management-api-service/internal/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func routes() *chi.Mux {
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
		r.Mount("/api/user", Routes.UserRoutes())
	})
	return router
}

func main() {
	configuration, err := config.GetEnv()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	router := routes()
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	//configuration type has two fields: Constants(PORT AND (URL + DBNAME)) & Database reference
	log.Println("Serving application at PORT :" + configuration.Constants.PORT)
	log.Fatal(http.ListenAndServe(":"+configuration.Constants.PORT, router))

}
