package services

import (
	"errors"
	"exercise-login-back-go/internal/model"
	"fmt"
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo model.UserRepository
}

func NewUserService(repo model.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(req model.UserRegistrationRequest) error {
	var err error
	if err := s.ValidateRegistration(req); err != nil {
		log.Println(err.Error())
		return err
	}
	// Crear el usuario y guardarlo en la base de datos
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	user.Password, err = hashPassword(req.Password)
	if err != nil {
		return err
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		log.Println(err.Error())
		return errors.New("error al crear el usuario")
	}
	return nil
}

func (s *UserService) ValidateRegistration(user model.UserRegistrationRequest) error {

	if !isValidEmail(user.Email) {
		return fmt.Errorf("el formato del correo electrónico no es válido")
	}
	if !isValidPhone(user.Phone) {
		return fmt.Errorf("el teléfono debe tener 10 dígitos")
	}
	if err := isValidPassword(user.Password); err != nil {
		return err
	}
	existingUser, err := s.repo.GetUserByEmailOrPhone(user.Email, user.Phone)
	if err != nil {
		return fmt.Errorf("error al verificar la existencia del usuario: %v", err)
	}
	if existingUser != nil {
		return errors.New("el correo/telefono ya se encuentra registrado")
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+`)
	return re.MatchString(email)
}

func isValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\d{10}$`)
	return re.MatchString(phone)
}

func isValidPassword(password string) error {
	var (
		hasMinLen  = len(password) >= 6
		hasMaxLen  = len(password) <= 12
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[@$&]`).MatchString(password)
	)

	if !hasMinLen {
		return fmt.Errorf("la contraseña debe tener al menos 6 caracteres")
	}
	if !hasMaxLen {
		return fmt.Errorf("la contraseña debe tener máximo 12 caracteres")
	}
	if !hasUpper {
		return fmt.Errorf("la contraseña debe incluir al menos una letra mayúscula")
	}
	if !hasLower {
		return fmt.Errorf("la contraseña debe incluir al menos una letra minúscula")
	}
	if !hasNumber {
		return fmt.Errorf("la contraseña debe incluir al menos un número")
	}
	if !hasSpecial {
		return fmt.Errorf("la contraseña debe incluir al menos un carácter especial")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	// La función GenerateFromPassword toma como argumentos la contraseña en forma de slice de bytes
	// y el costo del algoritmo de hashing. El costo es un factor de dificultad que define cuán segura
	// y "costosa" computacionalmente es la generación del hash. Un valor más alto hace que el hash sea
	// más seguro, pero también más lento de generar. El valor recomendado es 10.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err // Devuelve el error si la generación del hash falla
	}
	return string(hash), nil // Devuelve el hash como una cadena de texto
}

func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // Devuelve true si la contraseña coincide con el hash, false en caso contrario
}
