-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at,
    expires_at,
    revoked_at,
    user_id
) VALUES (
    $1,
    NOW(),
    NOW(),
    NOW() + INTERVAL '60 days',
    NULL,
    $2
)
RETURNING *;


-- name: GetUserFromRefreshToken :one
SELECT sqlc.embed(users), sqlc.embed(refresh_tokens) FROM refresh_tokens
JOIN users 
    ON users.id = refresh_tokens.user_id
WHERE token = $1 AND refresh_tokens.revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1; 