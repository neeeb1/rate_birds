-- name: CreateBird :one
INSERT INTO birds  (id, created_at, updated_at, common_name, scientific_name, family, "order", status)
VALUES ( 
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;
