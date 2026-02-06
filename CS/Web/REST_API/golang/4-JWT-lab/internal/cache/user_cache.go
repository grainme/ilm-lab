package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/movie-api/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	UserId   uuid.UUID
	Username string
	Role     domain.Role
}

// id: is the UUID of the refresh token
func GetUserByRefreshToken(ctx context.Context, rdb *redis.Client, refreshToken uuid.UUID) (*UserCache, error) {
	key := RefreshTokenKey(refreshToken)
	val, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		// cache miss (it's not an error)
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var userCache UserCache
	if err := json.Unmarshal([]byte(val), &userCache); err != nil {
		return nil, err
	}

	return &userCache, nil
}

// we store refresh token in Redis in case of "POST /auth/refresh"
// I'm thinking about do I need to save in Redis?
// should I save the ID and then get UserById?
// (but then it's expensive again because it requires a DB request)
// ---
// what if I store a some struct that contains (userId, role)
// the stuff, I need to generate another access token?
func SetUserByRefreshToken(ctx context.Context, rdb *redis.Client, refreshToken uuid.UUID, userCache UserCache) error {
	refreshTokenKey := RefreshTokenKey(refreshToken)

	data, err := json.Marshal(userCache)
	if err != nil {
		return err
	}

	// refresh token is valid for 7 days
	return rdb.Set(ctx, refreshTokenKey, data, time.Hour*24*7).Err()
}

func DelUserByRefreshTokenId(ctx context.Context, rdb *redis.Client, refreshTokenId uuid.UUID) error {
	refreshTokenKey := RefreshTokenKey(refreshTokenId)
	return rdb.Del(ctx, refreshTokenKey).Err()
}
