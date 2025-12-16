-- name: CreateBird :one
INSERT INTO birds  (id, created_at, updated_at, common_name, scientific_name, family, "order", status, image_urls)
VALUES ( 
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
) ON CONFLICT (common_name) DO UPDATE SET updated_at = NOW(), status = $5, image_urls = $6
RETURNING *;

-- name: GetRandomBird :many
SELECT * from birds
ORDER by RANDOM()
LIMIT $1;

-- name: GetRandomBirdWithImage :many
SELECT * from birds
WHERE image_urls is NOT NULL
ORDER by RANDOM()
LIMIT $1;

-- name: GetAllBirds :many
SELECT * from birds;

-- name: GetBirdByID :one
SELECT * from birds
WHERE id = $1;

-- name: GetTotalBirdCount :many
SELECT count(*) from birds;