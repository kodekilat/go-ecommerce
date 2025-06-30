package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodekilat/go-ecommerce/internal/models" // Ganti dengan path modul Anda
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := r.DB.QueryRow(context.Background(), query, user.FullName, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, full_name FROM users WHERE email = $1`
	var user models.User
	err := r.DB.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName)
	if err != nil {
		// Kita akan menangani error 'no rows' secara spesifik di handler
		return nil, err
	}
	return &user, nil
}
