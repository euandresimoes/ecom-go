-- +goose Up
-- +goose StatementBegin
ALTER TABLE product
ADD CONSTRAINT product_public_id_unique UNIQUE (public_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product
DROP CONSTRAINT product_public_id_unique;
-- +goose StatementEnd
