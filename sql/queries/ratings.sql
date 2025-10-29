-- name: PopulateRating :exec
INSERT INTO ratings (
    id,
    created_at,
    updated_at,
    matches,
    rating,
    bird_id
) VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
) ON CONFLICT (bird_id) DO NOTHING;