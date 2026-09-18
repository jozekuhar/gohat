-- name: ListOrders :many
SELECT * 
FROM orders
WHERE organization_id = @organization_id;
