package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/turanberker/tennis-league-service/internal/platform"
)

var redisCtx = context.Background()

func NewRedis() (*redis.Client, error) {
	config := platform.LoadRedisConfig()
	addr := fmt.Sprintf(
		"%s:%s",
		config.Host,
		config.Port,
	)
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

	// bağlantıyı doğrula
	if err := rdb.Ping(redisCtx).Err(); err != nil {
		return nil, err
	}

	fmt.Println("Redis connected")
	return rdb, nil
}
