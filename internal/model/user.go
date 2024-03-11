package model

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       int
	Username string
	Email    string
	Phone    string
	Password string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type UserRegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	EmailOrUsername string `json:"emailOrUsername"`
	Password        string `json:"password"`
}
