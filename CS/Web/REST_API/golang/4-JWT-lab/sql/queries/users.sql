-- name: AddUser :one
INSERT INTO
  users (id, username, password_hash)
VALUES
  ($1, $2, $3) RETURNING *;

-- name: FindUserByName :one
SELECT
  *
FROM
  users
WHERE
  username = $1;
