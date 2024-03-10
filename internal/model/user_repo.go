package model

type UserRepository interface {
	CreateUser(user User) error
	GetUserByEmailOrUsername(emailOrUsername string) (*User, error)
}
