package internal

import (
	"fmt"
	"math/rand"

	"github.com/thecsw/mira"
)

type RedditService struct {
	bot   *mira.Reddit
	cache *RedisCache
}

type redditImageEntry struct {
	author    string
	imageLink string
	width     int
	height    int
}

func (r RedditService) GetImage(animal string, width, height int) (ImageObject, error) {
	// todo, add getting images from reddit
	images, err := r.cache.GetImages(fmt.Sprintf("%s_%xx%x", animal, width, height))
	if err == nil && len(images) != 0 {
		randomIndex := rand.Intn(len(images))
		return images[randomIndex], nil
	}

	image, err := r.fetchImage(animal, width, height)
	if err == nil {
		r.cache.SaveImage(image.name, image.data)
		return image, nil
	}

	return ImageObject{}, err
}

func (r RedditService) fetchImage(animal string, width, height int) (ImageObject, error) {
	return ImageObject{}, nil
}

func (r RedditService) processImageEntry(entry redditImageEntry) []byte {
	return []byte{}
}

func NewRedditService(config *Config, cache *RedisCache) (*RedditService, error) {
	r, err := mira.Init(mira.Credentials{
		ClientId:     config.RedditAppID,
		ClientSecret: config.RedditSecret,
		Username:     config.RedditUsername,
		Password:     config.RedisPassword,
		UserAgent:    "github.com/lorow/placeitgo/:v0.1.0 (by /u/PlaceItBot)",
	})

	if err != nil {
		panic(fmt.Sprintf("something went wrong while starting reddit bot up - %s", err))
	}

	return &RedditService{bot: r}, nil
}
