package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // "localhost:6379"
	})

	// testing connection
	connTest := rdb.Ping(context.Background()).Err()
	return rdb, connTest
}
