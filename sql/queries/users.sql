-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password)
VALUES (gen_random_uuid(), $1, $2)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password = $3 WHERE id = $1
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpgradeUser :one
UPDATE users SET is_chirpy_red = TRUE, updated_at = NOW() WHERE id = $1
RETURNING *;
