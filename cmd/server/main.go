package main

import (
	"exercise-login-back-go/internal/api"
	"exercise-login-back-go/internal/config"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}
	// Create a Gin router
	r := mux.NewRouter()

	// Set up routes passing the database connection
	err = api.SetupRoutes(&cfg, r)
	if err != nil {
		log.Fatal("Error setting up routes: ", err)
	}
	// Define allowed headers, methods, and origins for CORS responses
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"*"}) // Adjust this to be more restrictive if necessary

	// Wrap the router with CORS middleware
	corsRouter := handlers.CORS(originsOk, headersOk, methodsOk)(r)

	// Start the HTTP server
	port := ":80"
	fmt.Println("Web server started on port", port)
	if err := http.ListenAndServe(port, corsRouter); err != nil {
		log.Println(err)
	}
}
