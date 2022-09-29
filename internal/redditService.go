package internal

import (
	"fmt"

	"github.com/thecsw/mira"
)

type RedditService struct {
	bot *mira.Reddit
}

type redditImageEntry struct {
	author    string
	imageLink string
	width     int
	height    int
}

func (r RedditService) GetImage(animal string, width, height int) (RedditImage, error) {
	return RedditImage{}, nil
}

func (r RedditService) processImageEntry(entry redditImageEntry) []byte {
	return []byte{}
}

func NewRedditService(config *Config) (*RedditService, error) {
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
