package api

import (
	"exercise-login-back-go/internal/config"
	"exercise-login-back-go/internal/repositories"
	"exercise-login-back-go/internal/services"
	"exercise-login-back-go/pkg/db"

	"github.com/gorilla/mux"
)

func SetupRoutes(cfg *config.Config, r *mux.Router) error {
	return configureLoginRoutes(cfg, r)
}

// configureLoginRoutes sets up the routes for the login service.
func configureLoginRoutes(cfg *config.Config, r *mux.Router) error {
	database, err := db.InitializeDatabase(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		return err
	}
	userRepository := repositories.NewUserRepository(database)

	// userService is the service used to handle user operations.
	userService := services.NewUserService(userRepository)

	// userHandler is the handler used to handle user requests.
	userHandler := NewUserHandler(userService)

	// Register the user registration handler.
	r.HandleFunc("/api/users/register", userHandler.RegisterUser).Methods("POST")

	return nil
}
