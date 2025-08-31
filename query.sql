-- name: AddInventory :one
INSERT INTO inventory (product_sku, stock_level)
VALUES ($1, $2)
ON CONFLICT (product_sku)
DO UPDATE SET stock_level = inventory.stock_level + EXCLUDED.stock_level
RETURNING product_sku, stock_level;

-- name: RemoveInventory :one
UPDATE inventory
SET stock_level = stock_level - $2
WHERE product_sku = $1
RETURNING product_sku, stock_level;

-- name: GetInventory :one
SELECT product_sku, stock_level
FROM inventory
WHERE product_sku = $1;
