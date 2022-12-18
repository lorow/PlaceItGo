package reddit

import (
	"fmt"
	"placeitgo/config"
	"placeitgo/model"
	"placeitgo/storage"

	"github.com/thecsw/mira"
)

type RedditService struct {
	bot   *mira.Reddit
	cache *storage.RedisCache
}

// new idea behind this, this just downloads the pic data and author info from reddit, processing is done in processor. Coordinating it all in an ImageService

func (r RedditService) GetImage(animal string, width, height int) (model.ImageResponse, error) {
	imageEntry, err := r.cache.GetImage(width, height, animal)

	if err == nil {
		// we got the iamge from cache, let's download it, process it, and return
		freshImage, err := r.fetchImage(imageEntry)
		if err != nil {
			return model.ImageResponse{},
				err
		}
		processedImage := r.processImageEntry(imageEntry.Author, imageEntry.Title, freshImage)
		return model.ImageResponse{Data: processedImage}, nil
	} else {
		// something went wrong with getting the image from cache, possibly nothing found
		// get it straight from reddit
		return model.ImageResponse{}, nil
	}
}

func (r RedditService) fetchImage(entry model.ImageDBEntry) ([]byte, error) {
	return []byte{}, nil
}

func (r RedditService) processImageEntry(author, title string, data []byte) []byte {
	return []byte{}
}

func NewRedditService(config *config.Config, cache *storage.RedisCache) (*RedditService, error) {
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
