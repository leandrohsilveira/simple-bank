
-- name: CreateUser :one
INSERT INTO users
    (email, password, name, role)
VALUES
    ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    name = $2,
    email = $3,
    role = $4
WHERE
    id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET
    password = $2
WHERE
    id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

