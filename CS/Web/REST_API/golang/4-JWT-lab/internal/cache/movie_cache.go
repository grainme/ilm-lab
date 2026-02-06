package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
	"github.com/redis/go-redis/v9"
)

func GetMovieById(ctx context.Context, rdb *redis.Client, id uuid.UUID) (*domain.Movie, error) {
	key := MovieKey(id)
	val, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		// cache miss (it's not an error)
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var movie domain.Movie
	if err := json.Unmarshal([]byte(val), &movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

func SetMovie(ctx context.Context, rdb *redis.Client, movieId uuid.UUID, movieVal domain.Movie) error {
	movieKey := MovieKey(movieId)

	data, err := json.Marshal(movieVal)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, movieKey, data, time.Minute*10).Err()
}

func DelMovie(ctx context.Context, rdb *redis.Client, movieId uuid.UUID) error {
	movieKey := MovieKey(movieId)
	return rdb.Del(ctx, movieKey).Err()
}

// atomic increment
func IncrementViewCount(ctx context.Context, rdb *redis.Client, id uuid.UUID) (int64, error) {
	key := ViewsMovieKey(id)
	return rdb.Incr(ctx, key).Result()
}
