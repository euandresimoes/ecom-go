package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db         *pgxpool.Pool
	redis      *redis.Client
	jwtManager *JWTManager
}

func NewRepository(db *pgxpool.Pool, redis *redis.Client, jwtManager *JWTManager) *Repository {
	return &Repository{db: db, redis: redis, jwtManager: jwtManager}
}

func (r *Repository) Register(data UserRegisterModel) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		return err
	}

	var exists bool
	query := `
		SELECT EXISTS
		(SELECT 1 FROM users WHERE email = $1)
	`
	err = r.db.QueryRow(
		context.Background(),
		query,
		data.Email,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("email already in use")
	}

	query = `
			INSERT INTO users (first_name, last_name, email, password_hash)
			VALUES ($1, $2, $3, $4)
		`
	_, err = r.db.Exec(
		context.Background(),
		query,
		data.FirstName,
		data.LastName,
		data.Email,
		string(hashedPassword),
	)

	return err
}

func (r *Repository) Login(data UserLoginModel) (string, error) {
	var (
		id            int
		role          UserRole
		password_hash string
	)

	query := `
		SELECT id, role, password_hash
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		data.Email,
	).Scan(
		&id,
		&role,
		&password_hash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("account not found")
		}

		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(data.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", errors.New("invalid credentials")
	}

	return r.jwtManager.Sign(id, role)
}
