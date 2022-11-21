package internal

import (
	"context"
	"fmt"
	"hash/fnv"
	"math/rand"
	"strconv"

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

func (r RedisCache) fetchImageEntryData(keys []string) (ImageEntry, error) {
	keysLen := len(keys)
	if keysLen == 0 {
		return ImageEntry{}, fmt.Errorf("couldn't find anything in cache")
	}

	randomKey := keys[rand.Intn(keysLen)]
	imageData, err := r.db.HGetAll(ctx, randomKey).Result()
	if err != nil {
		return ImageEntry{}, err
	}

	imageWidth, imageHeight, err := ConvertResolutionFormString(imageData["dimensions"], "x")

	if err != nil {
		return ImageEntry{}, err
	}

	return ImageEntry{
		author: imageData["author"],
		title:  imageData["title"],
		link:   imageData["link"],
		width:  imageWidth,
		height: imageHeight,
	}, nil
}

func (r RedisCache) GetImage(width, height int, animal string) (ImageEntry, error) {

	// todo, refactor it to be cleaner a bit
	var cursor uint64
	keys, _, err := r.db.Scan(ctx, cursor, fmt.Sprintf("image:*:%dx%d-%s*", width, height, animal), 10).Result()
	if err != nil {
		return ImageEntry{}, err
	}

	if len(keys) > 0 {
		return r.fetchImageEntryData(keys)
	} else {
		storedResolutions, err := r.db.SMembers(ctx, fmt.Sprintf("%s_resolutions", animal)).Result()
		if err != nil {
			return ImageEntry{}, err
		}

		matchingResolutions := []string{}
		for _, storedResolution := range storedResolutions {
			imageWidth, imageHeight, err := ConvertResolutionFormString(storedResolution, "x")

			if err != nil {
				return ImageEntry{}, err
			}
			// we're looking for "similar"resolutions, as in bigger than what was asked for
			// because we will be cropping them to size
			if (imageWidth-width) >= 0 && (imageHeight-height) >= 0 && len(matchingResolutions) < 3 {
				matchingResolutions = append(matchingResolutions, fmt.Sprintf("%d_%d", width, height))
			}
		}
		return r.fetchImageEntryData(matchingResolutions)
	}
}

func (r RedisCache) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (r RedisCache) SaveImage(width, height int, authorName, title, animal, imageLink string) error {
	linkHash := strconv.Itoa(int(r.hash(imageLink)))
	// the key looks like image:<link_hash>:1920x1080_fox_<detail>
	_, err := r.db.HSet(ctx, fmt.Sprintf("image:%s:%dx%d-%s", linkHash, width, height, animal), map[string]interface{}{
		"author": authorName,
		"title":  title,
		"link":   imageLink,
		"width":  width,
		"height": height,
	}).Result()

	if err != nil {
		return err
	}

	_, err = r.db.SAdd(ctx, fmt.Sprintf("%s_resolutions", animal), fmt.Sprintf("%dx%d", width, height)).Result()
	if err != nil {
		return err
	}
	return nil
}

func testRedisConnection(rdb *redis.Client) error {
	err := rdb.Set(ctx, "testConnection", 1, 0).Err()

	if err == nil {
		rdb.Del(ctx, "testConnection").Err()
	}

	return err
}

func NewRedisCache(cfg *Config) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisURL, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	})

	err := testRedisConnection(rdb)

	redisCache := RedisCache{
		db: rdb,
	}

	return &redisCache, err
}
