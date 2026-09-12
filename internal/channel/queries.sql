-- name: CreateChannel :one
INSERT INTO channels (id, organization_id, name, provider, credentials, status)
VALUES (@id, @organization_id, @name, @provider, @credentials, @status)
RETURNING *;
