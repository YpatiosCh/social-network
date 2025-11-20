-- name: GetUserProfile :one
SELECT
    id,
    username,
    first_name,
    last_name,
    date_of_birth,
    avatar,
    about_me,
    profile_public
FROM users
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetUserBasic :one
SELECT
    id,
    username,
    avatar,
    profile_public
FROM users
WHERE id = $1;
  
  
-- name: UpdateUserProfile :one
UPDATE users
SET
    username      = COALESCE($2, username),
    first_name    = COALESCE($3, first_name),
    last_name     = COALESCE($4, last_name),
    date_of_birth = COALESCE($5, date_of_birth),
    avatar        = COALESCE($6, avatar),
    about_me      = COALESCE($7, about_me),
    updated_at    = CURRENT_TIMESTAMP
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: UpdateProfilePrivacy :exec
UPDATE users
SET profile_public=$2
WHERE id=$1;

-- name: UpdateUserPassword :exec
UPDATE auth_user
SET
    password_hash = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateUserEmail :exec
UPDATE auth_user
SET
    email = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: SearchUsers :many
SELECT
    id,
    username,
    avatar,
    profile_public
FROM users
WHERE deleted_at IS NULL
  AND (
        username % $1 OR
        first_name % $1 OR
        last_name % $1
      )
ORDER BY similarity(username, $1) DESC,
         similarity(first_name, $1) DESC,
         similarity(last_name, $1) DESC
LIMIT $2;