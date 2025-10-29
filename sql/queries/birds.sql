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
) ON CONFLICT (common_name) DO UPDATE SET updated_at = NOW(), status = $5
RETURNING *;

-- name: GetRandomBird :many
SELECT * from birds
ORDER by RANDOM()
LIMIT $1;

-- name: GetAllBirds :many
SELECT * from birds;