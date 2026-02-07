package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisCtx = context.Background()

func NewRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // varsa doldur
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	// bağlantıyı doğrula
	if err := rdb.Ping(redisCtx).Err(); err != nil {
		return nil, err
	}

	fmt.Println("Redis connected")
	return rdb, nil
}