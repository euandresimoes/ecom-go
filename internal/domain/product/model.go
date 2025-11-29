package product

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductWeightUnit string

const (
	ProductUnitG  ProductWeightUnit = "g"
	ProductUnitKG ProductWeightUnit = "kg"
)

type ProductModel struct {
	ID          pgtype.Int4        `json:"id"`
	PublicID    string             `json:"public_id"`
	Name        string             `json:"name"`
	Price       pgtype.Numeric     `json:"price"`
	Stock       int                `json:"stock"`
	CategoryID  int                `json:"category_id"`
	WeightUnit  ProductWeightUnit  `json:"weight_unit"`
	WeightValue pgtype.Numeric     `json:"weight_value"`
	Images      []string           `json:"images"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type ProductCreateDto struct {
	Name        string            `json:"name" db:"name" validate:"required,min=3,max=30"`
	Price       pgtype.Numeric    `json:"price" db:"price" validate:"required"`
	Stock       int               `json:"stock" db:"stock" validate:"required"`
	CategoryID  int               `json:"category_id" db:"category_id" validate:"required"`
	WeightUnit  ProductWeightUnit `json:"weight_unit" db:"weight_unit" validate:"required"`
	WeightValue pgtype.Numeric    `json:"weight_value" db:"weight_value" validate:"required"`
	Images      []string          `json:"images" db:"images" validate:"required"`
}

type ProductUpdateDto struct {
	Name        *string            `json:"name" db:"name" validate:"omitempty,min=3,max=30"`
	Price       *pgtype.Numeric    `json:"price" db:"price" validate:"omitempty"`
	Stock       *int               `json:"stock" db:"stock" validate:"omitempty"`
	CategoryID  *int               `json:"category_id" db:"category_id" validate:"omitempty"`
	WeightUnit  *ProductWeightUnit `json:"weight_unit" db:"weight_unit" validate:"omitempty"`
	WeightValue *pgtype.Numeric    `json:"weight_value" db:"weight_value" validate:"omitempty"`
	Images      *[]string          `json:"images" db:"images" validate:"omitempty"`
}

type CategoryModel struct {
	ID   pgtype.Int4 `json:"id" db:"id"`
	Name string      `json:"name" db:"name"`
}

type CategoryCreateDto struct {
	Name string `json:"name" db:"name" validate:"required,min=3,max=50"`
}
