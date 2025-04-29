-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password)
VALUES (gen_random_uuid(), $1, $2)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password = $3 WHERE id = $1
RETURNING id, email, created_at, updated_at;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
