-- name: MakeRefreshToken :one
INSERT INTO refresh_tokens (
    token, created_at, updated_at, user_id, expires_at, revoked_at
)
VALUES (
    $1, $2, $3, $4, $5, CAST(NULL AS TIMESTAMP WITH TIME ZONE)
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;
