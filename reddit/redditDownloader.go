package reddit

import (
	"fmt"
	"placeitgo/config"
	"placeitgo/model"

	"github.com/thecsw/mira"
)

type RedditService struct {
	bot *mira.Reddit
}

func (r RedditService) GetImage(animal string, width, height int) (model.ImageDBEntry, error) {
	// // todo, add getting images from reddit
	return model.ImageDBEntry{}, nil

}

func NewRedditDownloader(config *config.Config) (*RedditService, error) {
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
