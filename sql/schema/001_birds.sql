-- +goose Up
CREATE TABLE birds (
  id UUID PRIMARY KEY,
  created_at TIMEsTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  common_name TEXT,
  scientific_name TEXT,
  family TEXT,
  "order" TEXT,
  status TEXT
);

-- +goose Down
DROP TABLE birds;
