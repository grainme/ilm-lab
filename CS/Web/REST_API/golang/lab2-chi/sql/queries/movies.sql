-- name: GetMovies :many
SELECT
  *
FROM
  movies;

-- name: GetMovieById :one
SELECT
  *
FROM
  movies
WHERE
  id = $1;

-- name: AddMovie :one
INSERT INTO
  movies (id, title, director, year)
VALUES
  ($1, $2, $3, $4) RETURNING *;

-- name: DeleteMovieById :exec
DELETE FROM movies
WHERE
  id = $1;
