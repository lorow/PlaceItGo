package storage

import (
	"context"
	"fmt"
	"hash/fnv"
	"math/rand"
	"placeitgo/config"
	"placeitgo/model"
	"placeitgo/utils"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Storage interface {
	GetImage(animal string, width, height int) (model.ImageDBEntry, error)
	SaveImageEntries(entries []model.ImageDBEntry, animal string) error
}

var ctx = context.Background()

type RedisCache struct {
	db *redis.Client
}

func (r RedisCache) fetchImageEntryData(keys []string) (model.ImageDBEntry, error) {
	keysLen := len(keys)
	if keysLen == 0 {
		return model.ImageDBEntry{}, fmt.Errorf("couldn't find anything in cache")
	}

	randomKey := keys[rand.Intn(keysLen)]
	imageData, err := r.db.HGetAll(ctx, randomKey).Result()
	if err != nil {
		return model.ImageDBEntry{}, err
	}

	imageWidth, imageHeight, err := utils.ConvertResolutionFormString(fmt.Sprintf("%sx%s", imageData["width"], imageData["height"]), "x")

	if err != nil {
		return model.ImageDBEntry{}, err
	}

	return model.ImageDBEntry{
		Author: imageData["author"],
		Title:  imageData["title"],
		Link:   imageData["link"],
		Width:  imageWidth,
		Height: imageHeight,
	}, nil
}

func (r RedisCache) GetImage(animal string, width, height int) (model.ImageDBEntry, error) {

	// todo, refactor it to be cleaner a bit
	var cursor uint64
	keys, _, err := r.db.Scan(ctx, cursor, fmt.Sprintf("image:*:%dx%d-%s", width, height, animal), 10).Result()
	if err != nil {
		return model.ImageDBEntry{}, err
	}

	if len(keys) > 0 {
		return r.fetchImageEntryData(keys)
	} else {
		storedResolutions, err := r.db.SMembers(ctx, fmt.Sprintf("%s_resolutions", animal)).Result()
		if err != nil {
			return model.ImageDBEntry{}, err
		}

		matchingResolutions := []string{}
		for _, storedResolution := range storedResolutions {
			imageWidth, imageHeight, err := utils.ConvertResolutionFormString(storedResolution, "x")

			if err != nil {
				return model.ImageDBEntry{}, err
			}
			// we're looking for "similar" resolutions, as in bigger than what was asked for,
			// while still not exceeding a certain threshold, because we will be cropping them to size
			matchingWidth := bool(imageWidth-width > 0 && float32(imageWidth-width) <= float32(imageWidth)*0.3)
			matchingHeight := bool(imageHeight-height > 0 && float32(imageHeight-height) <= float32(imageHeight)*0.3)
			if matchingWidth && matchingHeight && len(matchingResolutions) < 3 {
				keys, _, _ = r.db.Scan(ctx, cursor, fmt.Sprintf("image:*:%dx%d-%s", imageWidth, imageHeight, animal), 10).Result()

				if len(keys) > 0 {
					matchingResolutions = append(matchingResolutions, keys[0])
				}
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

func (r RedisCache) SaveImageEntries(entries []model.ImageDBEntry, animal string) error {
	for _, imageEntry := range entries {
		linkHash := strconv.Itoa(int(r.hash(imageEntry.Link)))
		// the key looks like image:<link_hash>:1920x1080_fox_<detail>
		_, err := r.db.HMSet(ctx, fmt.Sprintf("image:%s:%dx%d-%s", linkHash, imageEntry.Width, imageEntry.Height, animal), map[string]interface{}{
			"author": imageEntry.Author,
			"title":  imageEntry.Title,
			"link":   imageEntry.Link,
			"width":  imageEntry.Width,
			"height": imageEntry.Height,
		}).Result()

		if err != nil {
			return err
		}

		_, err = r.db.SAdd(ctx, fmt.Sprintf("%s_resolutions", animal), fmt.Sprintf("%dx%d", imageEntry.Width, imageEntry.Height)).Result()
		if err != nil {
			return err
		}
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

func NewRedisCache(cfg *config.Config) (*RedisCache, error) {
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
