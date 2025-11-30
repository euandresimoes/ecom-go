package database

import (
	"context"
	"errors"

	"github.com/euandresimoes/ecom-go/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func NewPostgres(url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func CreateAdmin(adminPwd string, db *pgxpool.Pool) error {
	var count int

	query := `
		SELECT COUNT(*) FROM users WHERE role = 'admin'
	`
	err := db.QueryRow(
		context.Background(),
		query,
	).Scan(
		&count,
	)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	if adminPwd == "" {
		return errors.New("ADMIN_PASSWORD env not set")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPwd), 10)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO users
		(first_name, last_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = db.Exec(
		context.Background(),
		query,
		"John", "Doe", "admin@admin.com", hash, models.RoleAdmin,
	)

	return err
}
