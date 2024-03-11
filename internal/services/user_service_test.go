package services_test

import (
	"exercise-login-back-go/internal/mocks"
	"exercise-login-back-go/internal/model"
	"exercise-login-back-go/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRegistration(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := services.NewUserService(mockRepo, "dummySecret")

	t.Run("Invalid Email", func(t *testing.T) {
		req := model.UserRegistrationRequest{
			Email:    "invalidEmail",
			Phone:    "1234567890",
			Password: "Password@123",
		}
		err := service.ValidateRegistration(req)
		assert.Error(t, err)
		assert.Equal(t, "el formato del correo electrónico no es válido", err.Error())
	})

	t.Run("Invalid Phone", func(t *testing.T) {
		req := model.UserRegistrationRequest{
			Email:    "test@example.com",
			Phone:    "123",
			Password: "Password@123",
		}
		err := service.ValidateRegistration(req)
		assert.Error(t, err)
		assert.Equal(t, "el teléfono debe tener 10 dígitos", err.Error())
	})

	t.Run("Invalid Password", func(t *testing.T) {
		req := model.UserRegistrationRequest{
			Email:    "test@example.com",
			Phone:    "1234567890",
			Password: "pass",
		}
		err := service.ValidateRegistration(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "la contraseña debe tener al menos 6 caracteres")
	})

	t.Run("User Exists", func(t *testing.T) {
		mockRepo.On("GetUserByEmailOrPhone", "test@example.com", "1234567890").Return(&model.User{}, nil)

		req := model.UserRegistrationRequest{
			Email:    "test@example.com",
			Phone:    "1234567890",
			Password: "Password@123",
		}
		err := service.ValidateRegistration(req)
		assert.Error(t, err)
		assert.Equal(t, "el correo/telefono ya se encuentra registrado", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
