package internal

import (
	"context"
	"fmt"

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
	keys := ""

	for i := range keys {
		keys += fmt.Sprintf(" %s_%d", resolution, i)
	}

	return keys
}

func (r RedisCache) GetImages(resolution string) ([][]byte, error) {
	keys := r.getImageKeys(resolution)

	cachedImages, err := r.db.MGet(ctx, keys).Result()
	images := [][]byte{}

	for _, image := range cachedImages {
		images = append(images, []byte(fmt.Sprintf("%v", image)))
	}

	if err != nil {
		return nil, err
	}

	return images, err
}

func testRedisConnection(rdb *redis.Client) error {
	err := rdb.Set(ctx, "testConnection", 1, 0).Err()

	if err == nil {
		rdb.Del(ctx, "testConnection").Err()
	}

	return err
}

func GetRedisCache(cfg Config) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.redis_url,
		Password: cfg.redis_password,
		DB:       cfg.redis_database,
	})

	err := testRedisConnection(rdb)

	redisCache := RedisCache{
		db: rdb,
	}

	return &redisCache, err
}
