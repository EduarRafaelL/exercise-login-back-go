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
	query := "EXEC CreateUser @Username = @p1, @Email = @p2, @Phone = @p3, @Password = @p4"
	_, err := r.db.Exec(query,
		sql.Named("p1", user.Username),
		sql.Named("p2", user.Email),
		sql.Named("p3", user.Phone),
		sql.Named("p4", user.Password))
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetUserByEmailOrUsername(emailOrUsername string) (*model.User, error) {
	return nil, nil
}

func (r *userRepository) GetUserByEmailOrPhone(email, phone string) (*model.User, error) {
	var user model.User
	query := "EXEC GetUserByEmailOrPhone @Email = @p1, @Phone = @p2"
	row := r.db.QueryRow(query, sql.Named("p1", email), sql.Named("p2", phone))

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró el usuario, lo cual no es necesariamente un error
		}
		// Manejo de otros errores
		return nil, err
	}

	return &user, nil
}
