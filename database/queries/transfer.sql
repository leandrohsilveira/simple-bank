
-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id, to_account_id, amount, status
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetTransferById :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountTransfers :one
SELECT COUNT(*) FROM transfers;

-- name: UpdateTransfer :one
UPDATE transfers
SET amount = $2, status = $3
WHERE id = $1
RETURNING *;

-- name: SetTransferStatus :one
UPDATE transfers
SET status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;

