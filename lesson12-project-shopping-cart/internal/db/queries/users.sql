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

-- name: UpdateUser :one
UPDATE users
SET
    user_password = COALESCE(sqlc.narg(user_password), user_password),
    user_fullname = COALESCE(sqlc.narg(user_fullname), user_fullname),
    user_age = COALESCE(sqlc.narg(user_age), user_age),
    user_status = COALESCE(sqlc.narg(user_status), user_status),
    user_level = COALESCE(sqlc.narg(user_level), user_level)
WHERE
    user_uuid = sqlc.arg(user_uuid)::uuid
    AND user_deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :one
UPDATE users
SET
    user_deleted_at = now()
WHERE
    user_uuid = sqlc.arg(user_uuid)::uuid
    AND user_deleted_at IS NULL
RETURNING *;

-- name: RestoreUser :one
UPDATE users
SET
    user_deleted_at = NULL
WHERE
    user_uuid = sqlc.arg(user_uuid)::uuid
    AND user_deleted_at IS NOT NULL
RETURNING *;

-- name: TrashUser :one
DELETE FROM users
WHERE
    user_uuid = sqlc.arg(user_uuid)::uuid
    AND user_deleted_at IS NOT NULL
RETURNING *;