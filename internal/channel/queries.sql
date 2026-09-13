-- name: CreateChannel :one
INSERT INTO channels (id, organization_id, name, provider, credentials, status)
VALUES (@id, @organization_id, @name, @provider, @credentials, @status)
RETURNING *;

-- name: ListChannels :many
SELECT *
FROM channels
WHERE organization_id = @organization_id;
