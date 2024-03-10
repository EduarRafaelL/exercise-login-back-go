package model

type UserRepository interface {
	CreateUser(user User) error
	GetUserByEmailOrUsername(emailOrUsername string) (*User, error)
	GetUserByEmailOrPhone(email, Phone string) (*User, error)
}
