CREATE INDEX IF NOT EXISTS movie_review_index ON REVIEWS (movie_id);

CREATE INDEX IF NOT EXISTS review_index ON REVIEWS (rating);
