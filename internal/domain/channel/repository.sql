-- name: CreateChannel :one
INSERT INTO channels (id, organization_id, name, provider, credentials, status)
VALUES (@id, @organization_id, @name, @provider, @credentials, @status)
RETURNING *;


-- name: ListChannels :many
SELECT *
FROM channels
WHERE organization_id = @organization_id
ORDER BY LOWER(name) ASC;


-- name: GetChannel :one
SELECT *
FROM channels
WHERE organization_id = @organization_id
  AND id = @channel_id;


-- name: DeleteChannel :exec
DELETE FROM channels
WHERE organization_id = @organization_id
  AND id = @channel_id;
