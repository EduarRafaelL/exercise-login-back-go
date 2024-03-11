package services

import (
	"errors"
	"exercise-login-back-go/internal/model"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      model.UserRepository
	SecretKey string
}

func NewUserService(repo model.UserRepository, secretKey string) *UserService {
	return &UserService{
		repo:      repo,
		SecretKey: secretKey,
	}
}

// RegisterUser creates a new user in the system.
func (s *UserService) RegisterUser(req model.UserRegistrationRequest) error {
	var err error
	// Validate the user registration request
	if err := s.ValidateRegistration(req); err != nil {
		return err
	}

	// Create the user
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	// Hash the password
	user.Password, err = hashPassword(req.Password)
	if err != nil {
		return err
	}

	// Save the user
	if err := s.repo.CreateUser(user); err != nil {
		log.Println(err.Error())
		return errors.New("error al crear el usuario")
	}

	return nil
}

// ValidateRegistration validates the user registration request.
func (s *UserService) ValidateRegistration(req model.UserRegistrationRequest) error {
	if !isValidEmail(req.Email) {
		return fmt.Errorf("el formato del correo electrónico no es válido")
	}
	if !isValidPhone(req.Phone) {
		return fmt.Errorf("el teléfono debe tener 10 dígitos")
	}
	if err := isValidPassword(req.Password); err != nil {
		return err
	}
	existingUser, err := s.repo.GetUserByEmailOrPhone(req.Email, req.Phone)
	if err != nil {
		return fmt.Errorf("error al verificar la existencia del usuario: %v", err)
	}
	if existingUser != nil {
		return errors.New("el correo/telefono ya se encuentra registrado")
	}
	return nil
}

// isValidEmail determines if the email is in the correct format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+`)
	return re.MatchString(email)
}

// isValidPhone determines if the phone number is in the correct format
func isValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\d{10}$`)
	return re.MatchString(phone)
}

// Check if the password meets the requirements
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

// hashPassword takes a password as a string and returns a hashed version of it as a string.
// It uses the bcrypt library to hash the password.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// LoginUser authenticates a user using their email or username and password.
// If the credentials are valid, a JSON Web Token (JWT) is generated and returned.
func (s *UserService) LoginUser(emailOrUsername, password string) (string, error) {
	// Retrieve the user from the database based on the provided email or username.
	user, err := s.repo.GetUserByEmailOrUsername(emailOrUsername)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Compare the provided password with the hashed password in the database.
	if err := s.checkPassword(user.Password, password); err != nil {
		return "", err
	}

	// Generate a JSON Web Token (JWT) and return it to the caller.
	tokenString, err := s.createToken(user.Username)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Verify the provided password
func (s *UserService) checkPassword(hashedPassword, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		return errors.New("contraseña incorrecta")
	}
	return nil
}

// Create a new JSON Web Token (JWT)
func (s *UserService) createToken(username string) (string, error) {
	// Set the expiration time for the token to be 1 hour from now
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create the claims for the token
	claims := &model.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "LOGIN-EXERCISE-TOKEN",
		},
	}

	// Create the JWT with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the JWT with the secret key
	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
