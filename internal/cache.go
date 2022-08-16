package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	GetImage(imageName string) ([]byte, error)
	CacheImage(imageData []byte, imageName string) error
}

var ctx = context.Background()

type RedisCache struct {
	db *redis.Client
}

func (r RedisCache) GetImage(imageName string) ([]byte, error) {
	image, err := r.db.Get(ctx, imageName).Result()
	return []byte(image), err
}

func GetRedisCache(url, password string, db int) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       db,
	})

	redisCache := RedisCache{
		db: rdb,
	}

	return &redisCache, nil
}
