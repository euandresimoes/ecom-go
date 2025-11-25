package product

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductModel struct {
	ID          pgtype.Int4        `json:"id" db:"id"`
	PublicID    string             `json:"public_id" db:"public_id"`
	Name        string             `json:"name" db:"name"`
	Description pgtype.Text        `json:"description" db:"description"`
	Price       pgtype.Numeric     `json:"price" db:"price"`
	CreatedAt   pgtype.Timestamptz `json:"created_at" db:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at" db:"updated_at"`
}

type ProductCreateDto struct {
	Name        pgtype.Text    `json:"name" db:"name"`
	Description pgtype.Text    `json:"description" db:"description"`
	Price       pgtype.Numeric `json:"price" db:"price"`
}

type ProductUpdateDto struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Price       *pgtype.Numeric `json:"price"`
}
