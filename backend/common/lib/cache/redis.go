package cache

import (
	"context"
	"fmt"
	"time"

	"tennis-league/common/lib/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	config := config.LoadRedisConfig()

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     config.Password,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	fmt.Println("Redis connected")
	return rdb, nil
}
