package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/euandresimoes/ecom-go/backend/internal/infra/cache"
	"github.com/euandresimoes/ecom-go/backend/internal/models"
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

func (r *Repository) CreateCategory(data *models.CategoryCreateDto) (models.CategoryModel, error) {
	var c models.CategoryModel

	query := `
		INSERT INTO
		categories (name)
		VALUES ($1)
		RETURNING id, name
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		data.Name,
	).Scan(
		&c.ID,
		&c.Name,
	)
	if err != nil {
		return c, err
	}

	redisKey := "products:categories"
	cache.DeleteUnique(r.redis, redisKey)

	return c, nil
}

func (r *Repository) GetAllCategories() ([]models.CategoryModel, error) {
	var cList []models.CategoryModel

	redisKey := "products:categories"
	cachedCategories, err := cache.Get[[]models.CategoryModel](r.redis, redisKey)
	if cachedCategories != nil {
		return *cachedCategories, err
	}

	query := `
		SELECT
		id, name
		FROM categories
	`
	rows, err := r.db.Query(
		context.Background(),
		query,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("no categories found")
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.CategoryModel
		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}

		cList = append(cList, category)
	}

	if len(cList) == 0 {
		return nil, errors.New("no categories found")
	}

	cache.Set(r.redis, redisKey, &cList)

	return cList, nil
}

func (r *Repository) DeleteCategory(id int) (models.CategoryModel, error) {
	var c models.CategoryModel

	query := `
		DELETE FROM categories
		WHERE id = $1
		RETURNING id, name
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&c.ID,
		&c.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c, errors.New("category not found")
		}

		return c, err
	}

	redisKey := "products:categories"
	cache.DeleteUnique(r.redis, redisKey)

	return c, nil
}

func (r *Repository) Create(data *models.ProductCreateDto) (models.ProductModel, error) {
	var p models.ProductModel

	publicID := cuid.New()
	query := `
		INSERT INTO
		products (public_id, name, price, stock, category_id, weight_unit, weight_value, images)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, public_id, name, price, stock, category_id, weight_unit, weight_value, images, created_at, updated_at
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		publicID, data.Name, data.Price, data.Stock, data.CategoryID, data.WeightUnit, data.WeightValue, data.Images,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.WeightUnit,
		&p.WeightValue,
		&p.Images,
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

func (r *Repository) Delete(id int) (models.ProductModel, error) {
	var p models.ProductModel

	query := `
		DELETE FROM products 
		WHERE id = $1
		RETURNING id, public_id, name, price, stock, category_id, weight_unit, weight_value, images, created_at, updated_at
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.WeightUnit,
		&p.WeightValue,
		&p.Images,
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

func (r *Repository) Update(id int, data *models.ProductUpdateDto) (models.ProductModel, error) {
	var p models.ProductModel

	query := `
		UPDATE products
		SET
			name = COALESCE($2, name),
			price = COALESCE($3, price),
			stock = COALESCE($4, stock),
			category_id = COALESCE($5, category_id),
			weight_unit = COALESCE($6, weight_unit),
			weight_value = COALESCE($7, weight_value),
			images = COALESCE($8, images),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, public_id, name, price, stock, category_id, weight_unit, weight_value, images, created_at, updated_at
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		id, data.Name, data.Price, data.Stock, data.CategoryID, data.WeightUnit, data.WeightValue, data.Images,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.WeightUnit,
		&p.WeightValue,
		&p.Images,
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

func (r *Repository) GetAll() ([]models.ProductModel, error) {
	var products []models.ProductModel

	redisKey := "products:all"
	cachedProducts, _ := cache.Get[[]models.ProductModel](r.redis, redisKey)
	if cachedProducts != nil {
		return *cachedProducts, nil
	}

	query := `
		SELECT id, public_id, name, price, stock, category_id, weight_unit, weight_value, images, created_at, updated_at
		FROM products
	`
	rows, err := r.db.Query(
		context.Background(),
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.ProductModel
		err := rows.Scan(
			&p.ID,
			&p.PublicID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&p.WeightUnit,
			&p.WeightValue,
			&p.Images,
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

	err = cache.Set(r.redis, redisKey, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) GetByID(id int) (models.ProductModel, error) {
	var p models.ProductModel

	redisKey := fmt.Sprintf("products:id:%v", id)
	cachedProduct, _ := cache.Get[models.ProductModel](r.redis, redisKey)
	if cachedProduct != nil {
		return *cachedProduct, nil
	}

	query := `
		SELECT id, public_id, name, price, stock, category_id, weight_unit, weight_value, images, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.WeightUnit,
		&p.WeightValue,
		&p.Images,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, errors.New("product not found")
		}

		return p, err
	}

	err = cache.Set(r.redis, redisKey, &p)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (r *Repository) GetByPublicID(publicID string) (models.ProductModel, error) {
	var p models.ProductModel

	redisKey := fmt.Sprintf("products:public:%v", publicID)
	cachedProduct, _ := cache.Get[models.ProductModel](r.redis, redisKey)
	if cachedProduct != nil {
		return *cachedProduct, nil
	}

	query := `
		SELECT id, public_id, name, price, stock, category_id, weight_unit, weight_value, images, created_at, updated_at
		FROM products
		WHERE public_id = $1
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		publicID,
	).Scan(
		&p.ID,
		&p.PublicID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.WeightUnit,
		&p.WeightValue,
		&p.Images,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, errors.New("product not found")
		}

		return p, err
	}

	err = cache.Set(r.redis, redisKey, &p)
	if err != nil {
		return p, err
	}

	return p, nil
}
