-- +goose Up
ALTER TABLE ratings
ADD COLUMN common_name TEXT UNIQUE
    REFERENCES birds(common_name) 
    ON DELETE CASCADE;

UPDATE ratings
SET common_name = birds.common_name
FROM birds
WHERE ratings.bird_id = birds.id;

-- +goose Down
ALTER TABLE ratings
DROP COLUMN common_name;