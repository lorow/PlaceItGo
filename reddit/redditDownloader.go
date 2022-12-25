package reddit

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"placeitgo/config"
	"placeitgo/model"
	"placeitgo/utils"

	"github.com/thecsw/mira"
)

type RedditService struct {
	bot *mira.Reddit
}

func (r RedditService) GetImage(animal string, width, height int) (model.ImageDBEntry, error) {
	// todo, add getting images from reddit
	var desiredSubmision *mira.PostListingChild
	var actualImageWidth int
	var actualImageHeight int

	// todo add max tries
	submissions, err := r.bot.GetSubredditPosts(animal, "hot", "all", 10)

	if err != nil {
		return model.ImageDBEntry{}, err
	}

	for _, submission := range submissions {
		imageLink := submission.GetText()
		imageData, err := utils.FetchImageFromURL(imageLink)
		if err != nil {
			continue
		}

		downloadedImageConfig, _, err := image.DecodeConfig(bytes.NewReader(imageData))
		if err != nil {
			continue
		}

		if downloadedImageConfig.Width >= width && downloadedImageConfig.Height >= height {
			actualImageWidth = downloadedImageConfig.Width
			actualImageHeight = downloadedImageConfig.Height
			desiredSubmision = &submission
		}
	}

	return model.ImageDBEntry{
		Author: desiredSubmision.GetAuthor(),
		Title:  desiredSubmision.GetTitle(),
		Link:   desiredSubmision.GetText(),
		Width:  actualImageWidth,
		Height: actualImageHeight,
	}, nil

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
