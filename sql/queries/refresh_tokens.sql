-- name: SaveRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: GetUserIdFromValidRefreshToken :one
SELECT user_id FROM refresh_tokens
WHERE token = $1 AND revoked_at IS NULL AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1 AND user_id = $2 AND revoked_at IS NULL;
