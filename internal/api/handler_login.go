package api

import (
	"encoding/json"
	"exercise-login-back-go/internal/model"
	"exercise-login-back-go/internal/services"
	"net/http"
)

type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new instance of a user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterUser registers a new user
func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// req is the incoming user registration request
	var req model.UserRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}
	// Write a status code of 200 OK
	w.WriteHeader(http.StatusOK)
	// Write a response of "Usuario registrado exitosamente"
	w.Write([]byte("Usuario registrado exitosamente"))
}
