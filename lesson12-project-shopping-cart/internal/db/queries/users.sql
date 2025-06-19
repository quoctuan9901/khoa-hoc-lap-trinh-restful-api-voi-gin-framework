-- name: GetUser :one
SELECT * FROM users WHERE uuid = $1;

-- name: CreateUser :one
INSERT INTO users(name, email) VALUES ($1, $2) RETURNING *;