-- name: CreateUser :one
INSERT INTO users (
    user_email,
    user_password,
    user_fullname,
    user_age,
    user_status,
    user_level
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;