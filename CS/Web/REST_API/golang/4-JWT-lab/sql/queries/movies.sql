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

-- name: UpdateMovieTitleById :exec
UPDATE movies
SET
  title = $2
WHERE
  id = $1;

-- name: DeleteMovieById :exec
DELETE FROM movies
WHERE
  id = $1;

-- name: GetMovieWithReviews :one
SELECT
  m.*,
  AVG(r.rating) as avg_rating,
  COUNT(r.*) as reviews_count
FROM
  movies m
  LEFT JOIN reviews r ON m.id = r.movie_id
WHERE
  m.id = $1
GROUP BY
  m.id;
