-- name: GetAllReviews :many
SELECT
  *
FROM
  REVIEWS;

-- name: AddReview :one
INSERT INTO
  REVIEWS (id, user_name, rating, comment, movie_id)
VALUES
  ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetAllReviewsByMovieId :many
SELECT
  *
FROM
  REVIEWS
WHERE
  movie_id = $1;
