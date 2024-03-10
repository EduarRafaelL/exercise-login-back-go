package repositories

import (
	"database/sql"
	"exercise-login-back-go/internal/model"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates and returns a new instance of the user repository.
func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}
func (r *userRepository) CreateUser(user model.User) error {
	return nil
}

func (r *userRepository) GetUserByEmailOrUsername(emailOrUsername string) (*model.User, error) {
	return nil, nil
}
