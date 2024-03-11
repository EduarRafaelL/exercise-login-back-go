package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"exercise-login-back-go/internal/api"
	"exercise-login-back-go/internal/mocks"
	"exercise-login-back-go/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	mockUserService := new(mocks.UserService)
	handler := api.NewUserHandler(mockUserService)

	t.Run("success", func(t *testing.T) {
		mockUserService.On("RegisterUser", mock.AnythingOfType("model.UserRegistrationRequest")).Return(nil)

		body, _ := json.Marshal(model.UserRegistrationRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Phone:    "1234567890",
			Password: "Password!23",
		})
		req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		handler.RegisterUser(resp, req)

		responseBody := resp.Body.String()
		t.Log("Success response body:", responseBody)
		assert.Equal(t, http.StatusOK, resp.Code)
		mockUserService.AssertExpectations(t)
	})

	t.Run("error decoding", func(t *testing.T) {
		body := []byte(`{bad json}`)
		req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		handler.RegisterUser(resp, req)

		responseBody := resp.Body.String()
		t.Log("Error decoding response body:", responseBody)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("validation errors", func(t *testing.T) {
		mockUserService.On("RegisterUser", mock.AnythingOfType("model.UserRegistrationRequest")).Return(errors.New("validation failed"))

		body, _ := json.Marshal(model.UserRegistrationRequest{
			Username: "testuser",
			Email:    "",
			Phone:    "1234567890",
			Password: "Password!23",
		})
		req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		handler.RegisterUser(resp, req)

		responseBody := resp.Body.String()
		t.Log("Validation error response body:", responseBody)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestLoginUser(t *testing.T) {
	mockUserService := new(mocks.UserService)
	handler := api.NewUserHandler(mockUserService)

	t.Run("login success", func(t *testing.T) {
		expectedToken := "fakeToken123"
		mockUserService.On("LoginUser", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(expectedToken, nil)

		body, _ := json.Marshal(model.UserLoginRequest{
			EmailOrUsername: "test@example.com",
			Password:        "Password!23",
		})
		req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		handler.LoginUser(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), expectedToken)
	})

	t.Run("error decoding request body", func(t *testing.T) {
		body := []byte(`{bad json}`)
		req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		handler.LoginUser(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "Error al decodificar el cuerpo de la solicitud")
	})

	t.Run("missing email or username", func(t *testing.T) {
		body, _ := json.Marshal(model.UserLoginRequest{
			EmailOrUsername: "",
			Password:        "Password!23",
		})
		req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		handler.LoginUser(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "Falta el campo email o nombre de usuario")
	})

	t.Run("missing password", func(t *testing.T) {
		body, _ := json.Marshal(model.UserLoginRequest{
			EmailOrUsername: "test@example.com",
			Password:        "",
		})
		req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		handler.LoginUser(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "Falta el campo contraseña")
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockUserService.ExpectedCalls = nil
		mockUserService.Calls = nil

		mockUserService.On("LoginUser", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("", errors.New("usuario / contraseña incorrectos"))

		body, _ := json.Marshal(model.UserLoginRequest{
			EmailOrUsername: "wrong@example.com",
			Password:        "wrongPassword",
		})
		req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()

		handler.LoginUser(resp, req)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "usuario / contraseña incorrectos")
	})
}
