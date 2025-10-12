-- +goose Up
CREATE TABLE birds (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  common_name TEXT UNIQUE,
  scientific_name TEXT UNIQUE,
  family TEXT,
  "order" TEXT,
  status TEXT
);

-- +goose Down
DROP TABLE birds;
