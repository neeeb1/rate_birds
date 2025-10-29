-- +goose Up
CREATE TABLE ratings (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    matches INTEGER DEFAULT 0,
    rating INTEGER DEFAULT 1000,
    bird_id UUID UNIQUE NOT NULL REFERENCES birds
    ON DELETE CASCADE,
    FOREIGN KEY (bird_id)
    REFERENCES birds(id)
);

-- +goose Down
DROP TABLE ratings;