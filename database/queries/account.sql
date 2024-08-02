
-- name: CreateAccount :one
INSERT INTO accounts (
    owner_id, balance, currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccountById :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1 
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountAccounts :one
SELECT COUNT(*) FROM accounts;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + @Change
WHERE id = @Id
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;