package auth

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

type UserModel struct {
	ID           int                `json:"id"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"password_hash"`
	Role         UserRole           `json:"role"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}

type UserPublicModel struct {
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Role      UserRole           `json:"role"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type UserRegisterModel struct {
	FirstName string `json:"first_name" db:"first_name" validate:"required,min=3,max=15"`
	LastName  string `json:"last_name" db:"last_name" validate:"omitempty,min=3,max=15"`
	Email     string `json:"email" db:"email" validate:"required,email,min=5,max=50"`
	Password  string `json:"password" db:"password_hash" validate:"required,min=8,max=32"`
}

type UserLoginModel struct {
	Email    string `json:"email" db:"email" validate:"required,email,max=50"`
	Password string `json:"password" db:"password_hash" validate:"required,min=8,max=32"`
}

type UserUpdateModel struct {
	FirstName   *string `json:"first_name" db:"first_name" validate:"omitempty,min=3,max=15"`
	LastName    *string `json:"last_name" db:"last_name" validate:"omitempty,min=3,max=15"`
	Email       *string `json:"email" db:"email" validate:"omitempty,email,min=5,max=50"`
	NewPassword *string `json:"new_password" db:"password_hash" validate:"omitempty,min=8,max=32"`
}
