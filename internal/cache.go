package internal

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

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

func (r RedisCache) getImageKeys(resolution string) string {
	keys := make([]string, 5)

	for i := 0; i < len(keys); i++ {
		keys[i] = fmt.Sprintf("%s_%d", resolution, i)
	}

	return strings.Join(keys, " ")
}

func (r RedisCache) GetImages(resolution string) ([][]byte, error) {
	keys := r.getImageKeys(resolution)

	images, err := r.db.MGet(ctx, keys).Result()

	if err != nil {
		return nil, err
	}
	randomIndex := rand.Intn(len(images))
	image := fmt.Sprint(images[randomIndex])

	return [][]byte{[]byte(image)}, err
}

func GetRedisCache(cfg Config) *RedisCache {
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
