-- +goose Up
ALTER TABLE birds
ADD COLUMN image_urls TEXT[];

-- +goose Down
ALTER TABLE birds
DROP COLUMN image_urls;