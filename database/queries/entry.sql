
-- name: CreateEntry :one
INSERT INTO entries (
    account_id, amount
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetEntryById :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountEntries :one
SELECT COUNT(*) FROM entries;

-- name: UpdateEntry :one
UPDATE entries
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;