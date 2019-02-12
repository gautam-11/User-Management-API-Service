package main

import (
	"log"
	"net/http"

	"user-management-api-service/Routes"
	"user-management-api-service/internal/config"

	"github.com/go-chi/chi"
)

//Startup function
func main() {
	configuration, err := config.GetEnv()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	router := Routes.SetRoutes()
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Printf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	//configuration type has two fields: Constants(PORT AND (URL + DBNAME)) & Database reference
	log.Println("Serving application at PORT :" + configuration.Constants.PORT)
	log.Fatal(http.ListenAndServe(":"+configuration.Constants.PORT, router))

}
