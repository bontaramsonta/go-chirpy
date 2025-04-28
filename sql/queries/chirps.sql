-- name: CreateChirp :one
INSERT INTO chirps (user_id, body)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllChirps :many
SELECT id, user_id, body, created_at, updated_at FROM chirps
ORDER BY created_at ASC;
