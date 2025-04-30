-- name: CreateChirp :one
INSERT INTO chirps (user_id, body)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllChirps :many
SELECT id, user_id, body, created_at, updated_at FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsByAuthorID :many
SELECT id, user_id, body, created_at, updated_at FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT id, user_id, body, created_at, updated_at FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;
