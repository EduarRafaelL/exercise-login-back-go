package model

// User representa un usuario en el sistema.
type User struct {
	ID       int
	Username string
	Email    string
	Phone    string
	Password string
}

// UserLoginRequest estructura de los datos de registro de usuario.
type UserRegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// UserLoginRequest estructura de los datos de solicitud de login de usuario.
type UserLoginRequest struct {
	EmailOrUsername string `json:"emailOrUsername"`
	Password        string `json:"password"`
}
