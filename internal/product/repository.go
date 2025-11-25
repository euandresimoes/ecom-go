package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/euandresimoes/ecom-go/internal/cache"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucsky/cuid"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewRepository(db *pgxpool.Pool, redis *redis.Client) *Repository {
	return &Repository{db: db, redis: redis}
}

func (r *Repository) Create(data ProductCreateDto) (ProductModel, error) {
	var p ProductModel

	publicID := cuid.New()

	err := r.db.QueryRow(
		context.Background(),
		"INSERT INTO product (public_id, name, description, price) VALUES ($1, $2, $3, $4) RETURNING id, public_id, name, description, price, created_at, updated_at",
		publicID, data.Name, data.Description, data.Price,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return p, err
	}

	redisKey := "products:*"
	err = cache.DeleteMany(r.redis, redisKey)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (r *Repository) Delete(id int) (ProductModel, error) {
	var p ProductModel

	err := r.db.QueryRow(
		context.Background(),
		"DELETE FROM product WHERE id = $1 RETURNING id, public_id, name, description, price, created_at, updated_at",
		id,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, errors.New("product not found")
		}

		return p, err
	}

	redisKey := "products:*"
	cache.DeleteUnique(r.redis, redisKey)

	return p, nil
}

func (r *Repository) Update(id int, data ProductUpdateDto) (ProductModel, error) {
	var p ProductModel

	query := `
		UPDATE product
		SET
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			price = COALESCE($3, price),
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, public_id, name, description, price, created_at, updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		data.Name, data.Description, data.Price, id,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, errors.New("product not found")
		}
		return p, err
	}

	redisKey := "products:*"
	cache.DeleteUnique(r.redis, redisKey)

	return p, nil
}

func (r *Repository) GetAll() ([]ProductModel, error) {
	var products []ProductModel

	redisKey := "products:all"
	cachedProducts, _ := cache.Get[[]ProductModel](r.redis, redisKey)
	if cachedProducts != nil {
		return *cachedProducts, nil
	}

	rows, err := r.db.Query(
		context.Background(),
		"SELECT id, public_id, name, description, price, created_at, updated_at FROM product",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p ProductModel
		err := rows.Scan(
			&p.ID,
			&p.PublicID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	if len(products) == 0 {
		return nil, errors.New("no products found")
	}

	err = cache.Set(r.redis, redisKey, products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) GetByID(id int) (ProductModel, error) {
	var p ProductModel

	redisKey := fmt.Sprintf("products:id:%v", id)
	cachedProduct, _ := cache.Get[ProductModel](r.redis, redisKey)
	if cachedProduct != nil {
		return *cachedProduct, nil
	}

	err := r.db.QueryRow(
		context.Background(),
		"SELECT id, public_id, name, description, price, created_at, updated_at FROM product WHERE id = $1",
		id,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, errors.New("product not found")
		}

		return p, err
	}

	err = cache.Set(r.redis, redisKey, p)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (r *Repository) GetByPublicID(publicID string) (ProductModel, error) {
	var p ProductModel

	redisKey := fmt.Sprintf("products:public:%v", publicID)
	cachedProduct, _ := cache.Get[ProductModel](r.redis, redisKey)
	if cachedProduct != nil {
		return *cachedProduct, nil
	}

	err := r.db.QueryRow(
		context.Background(),
		"SELECT id, public_id, name, description, price, created_at, updated_at FROM product WHERE public_id = $1",
		publicID,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, errors.New("product not found")
		}

		return p, err
	}

	err = cache.Set(r.redis, redisKey, p)
	if err != nil {
		return p, err
	}

	return p, nil
}
