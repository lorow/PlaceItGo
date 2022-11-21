package internal

import (
	"fmt"

	"github.com/thecsw/mira"
)

type RedditService struct {
	bot   *mira.Reddit
	cache *RedisCache
}

type ImageEntry struct {
	author string
	title  string
	link   string
	width  int
	height int
}

func (r RedditService) GetImage(animal string, width, height int) ([]byte, error) {
	// todo, add getting images from reddit

	imageEntry, err := r.cache.GetImage(width, height, animal)

	if err == nil {
		// we got the iamge from cache, let's download it, process it, and return
		freshImage, err := r.fetchImage(imageEntry)
		if err != nil {
			return []byte{},
				err
		}
		processedImage := r.processImageEntry(imageEntry.author, imageEntry.title, freshImage)
		return processedImage, nil
	} else {
		// something went wrong with getting the image from cache, possibly nothing found
		// get it straight from reddit
		return []byte{}, nil
	}
}

func (r RedditService) fetchImage(entry ImageEntry) ([]byte, error) {
	return []byte{}, nil
}

func (r RedditService) processImageEntry(author, title string, data []byte) []byte {
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
