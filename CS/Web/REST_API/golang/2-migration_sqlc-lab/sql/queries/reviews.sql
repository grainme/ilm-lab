-- name: GetAllReviews :many
SELECT
  *
FROM
  reviews;

-- name: AddReview :one
INSERT INTO
  reviews (id, user_name, rating, comment, movie_id)
VALUES
  ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetAllReviewsByMovieId :many
SELECT
  *
FROM
  reviews
WHERE
  movie_id = $1;
