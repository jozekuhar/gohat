-- name: ListOrdersWithDetails :many
SELECT sqlc.embed(o), sqlc.embed(c)
FROM orders AS o
JOIN channels AS c ON o.channel_id = c.id
WHERE o.organization_id = @organization_id;
