package api

import (
	"encoding/json"
	"exercise-login-back-go/internal/model"
	"exercise-login-back-go/internal/services"
	"net/http"
	"strings"
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
		respondWithError(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	if errs := validateRegistrationRequest(req); len(errs) > 0 {
		respondWithMultipleErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err := uh.userService.RegisterUser(req); err != nil {
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	// Write a response of "Usuario registrado exitosamente"
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Usuario registrado exitosamente"})
}

// validateRegistrationRequest valida los campos necesarios de la solicitud de registro.
func validateRegistrationRequest(req model.UserRegistrationRequest) []string {
	var errs []string
	if strings.TrimSpace(req.Username) == "" {
		errs = append(errs, "Falta el campo username o está en blanco")
	}
	if strings.TrimSpace(req.Email) == "" {
		errs = append(errs, "Falta el campo email o está en blanco")
	}
	if strings.TrimSpace(req.Phone) == "" {
		errs = append(errs, "Falta el campo phone o está en blanco")
	}
	if strings.TrimSpace(req.Password) == "" {
		errs = append(errs, "Falta el campo password o está en blanco")
	}
	return errs
}

// respondWithError envía una respuesta de error con el mensaje proporcionado.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
}

// respondWithMultipleErrors envía una respuesta con múltiples mensajes de error.
func respondWithMultipleErrors(w http.ResponseWriter, code int, errs []string) {
	respondWithJSON(w, code, map[string][]string{"errors": errs})
}

// respondWithJSON envía una respuesta JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
