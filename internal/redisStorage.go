package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func (r RedisCache) getImageKeys(resolution string) ([]string, error) {
	var cursor uint64 // todo iterate over this
	keys, _, err := r.db.Scan(ctx, cursor, fmt.Sprintf("*%s*", resolution), 19).Result()

	return keys, err
}

func (r RedisCache) GetImages(resolution string) ([]RedditImage, error) {
	keys, err := r.getImageKeys(resolution)

	if err != nil {
		return []RedditImage{}, err
	}

	cachedImages, err := r.db.MGet(ctx, strings.Join(keys, " ")).Result()
	images := []RedditImage{}

	for index, image := range cachedImages {
		image_data, ok := image.([]byte)
		if ok {
			images = append(images, RedditImage{
				data: image_data,
				name: keys[index],
			})
		}
	}

	if err != nil {
		return nil, err
	}

	return images, err
}

func (r RedisCache) SaveImage(key string, data []byte) error {
	_, err := r.db.Set(ctx, key, data, 24*time.Hour).Result()
	return err
}

func testRedisConnection(rdb *redis.Client) error {
	err := rdb.Set(ctx, "testConnection", 1, 0).Err()

	if err == nil {
		rdb.Del(ctx, "testConnection").Err()
	}

	return err
}

func NewRedisCache(cfg Config) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	})

	err := testRedisConnection(rdb)

	redisCache := RedisCache{
		db: rdb,
	}

	return &redisCache, err
}
