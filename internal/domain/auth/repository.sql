-- name: CreateUser :one
INSERT INTO users (id, email)
VALUES (@id, @email)
RETURNING *;


-- name: CreateAuthentication :one
INSERT INTO authentications (id, user_id, provider, provider_id, password_hash)
VALUES (@id, @user_id, @provider, @provider_id, @password_hash)
RETURNING *;


-- name: GetAuthenticationByEmail :one
SELECT a.*
FROM authentications AS a
LEFT JOIN users AS u ON a.user_id = u.id
WHERE u.email = @email
  AND a.provider = @provider;


-- name: GetAuthenticationByProvider :one
SELECT *
FROM authentications
WHERE provider = @provider
  AND provider_id = @provider_id;


-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expires_at)
VALUES (@id, @user_id, @expires_at)
RETURNING *;


-- name: GetSession :one
SELECT *
FROM sessions
WHERE id = @session_id;


-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = @session_id;
