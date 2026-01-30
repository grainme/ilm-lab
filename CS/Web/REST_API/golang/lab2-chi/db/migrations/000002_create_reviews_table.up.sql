-- reviews table: id, movie_id (FOREIGN KEY), user_name, rating (1-10), comment
--
CREATE TABLE IF NOT EXISTS REVIEWS (
  id UUID PRIMARY KEY,
  user_name TEXT NOT NULL,
  rating INTEGER NOT NULL CHECK (
    rating >= 1
    AND rating <= 10
  ),
  comment TEXT,
  movie_id UUID NOT NULL REFERENCES MOVIES (id)
)
