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

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterUser registers a new user
func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// req is the incoming request body
	var req model.UserRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	// validate the incoming request
	if errs := validateRegistrationRequest(req); len(errs) > 0 {
		respondWithMultipleErrors(w, http.StatusBadRequest, errs)
		return
	}

	// register the user
	if err := uh.userService.RegisterUser(req); err != nil {
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}

	// return a successful response
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Usuario registrado exitosamente"})
}

// LoginUser handles user login requests
func (uh *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// req is the incoming request body
	var loginReq model.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding request body")
		return
	}

	// validate request
	if strings.TrimSpace(loginReq.EmailOrUsername) == "" {
		respondWithError(w, http.StatusBadRequest, "email or username is required")
		return
	}
	if strings.TrimSpace(loginReq.Password) == "" {
		respondWithError(w, http.StatusBadRequest, "password is required")
		return
	}

	// attempt to login user
	token, err := uh.userService.LoginUser(loginReq.EmailOrUsername, loginReq.Password)
	if err != nil {
		// handle error
		respondWithError(w, http.StatusUnauthorized, "username or password is incorrect")
		return
	}

	// return response
	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

// validateRegistrationRequest validates the incoming user registration request
func validateRegistrationRequest(req model.UserRegistrationRequest) []string {
	var errs []string
	if strings.TrimSpace(req.Username) == "" {
		errs = append(errs, "username is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		errs = append(errs, "email is required")
	}
	if strings.TrimSpace(req.Phone) == "" {
		errs = append(errs, "phone is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		errs = append(errs, "password is required")
	}
	return errs
}

// respondWithError sends a response with an error message in Json format
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
}

// respondWithMultipleErrors sends a response with multiple error messages in JSON format
func respondWithMultipleErrors(w http.ResponseWriter, code int, errs []string) {
	respondWithJSON(w, code, map[string][]string{"errors": errs})
}

// respondWithJSON sends a response with a JSON payload.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
