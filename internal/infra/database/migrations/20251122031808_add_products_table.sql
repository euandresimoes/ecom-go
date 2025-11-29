-- +goose Up
-- +goose StatementBegin
CREATE TYPE weight_unit
AS ENUM 
('g', 'kg');

CREATE TABLE
IF NOT EXISTS
categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE
IF NOT EXISTS
products (
    id SERIAL PRIMARY KEY,
    public_id VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(30) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    category_id INT REFERENCES categories(id) NOT NULL,
    weight_unit weight_unit NOT NULL,
    weight_value DECIMAL(10, 2) NOT NULL,
    images TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TYPE IF EXISTS weight_unit;
-- +goose StatementEnd
