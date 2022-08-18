package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	GetImages(resolution string) ([][]byte, error)
	CacheImage(imageData []byte, imageName string) error
}

var ctx = context.Background()

type RedisCache struct {
	db *redis.Client
}

func (r RedisCache) GetImages(resolution string) ([][]byte, error) {
	// todo implment searching for multiple keys matching 
	image, err := r.db.Get(ctx, resolution).Result()
	return [][]byte{[]byte(image)}, err
}

func GetRedisCache(cfg Config) (*RedisCache) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.redis_url,
		Password: cfg.redis_password,
		DB:       cfg.redis_database,
	})

	redisCache := RedisCache{
		db: rdb,
	}

	return &redisCache
}
