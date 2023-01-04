package reddit

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"placeitgo/config"
	"placeitgo/model"

	"github.com/thecsw/mira"
)

type RedditService struct {
	bot *mira.Reddit
}

func (r RedditService) GetImage(animal string, width, height int) (model.ImageDBEntry, error) {
	var desiredSubmision *mira.PostListingChild
	var actualImageWidth int
	var actualImageHeight int

	// todo add max tries
	submissions, err := r.bot.GetSubredditPosts(animal, "hot", "all", 10)

	if err != nil {
		return model.ImageDBEntry{}, err
	}

	for _, submission := range submissions {
		if !submission.Data.Preview.Enabled {
			continue
		}

		if int(submission.Data.Preview.Images[0].Source.Width) >= width && int(submission.Data.Preview.Images[0].Source.Height) >= height {
			actualImageWidth = int(submission.Data.Preview.Images[0].Source.Width)
			actualImageHeight = int(submission.Data.Preview.Images[0].Source.Height)
			desiredSubmision = &submission
			break
		}
	}

	if desiredSubmision == nil {
		return model.ImageDBEntry{}, fmt.Errorf("could not find desired placeholder iamge for %s %dx%d", animal, width, height)
	}

	return model.ImageDBEntry{
		Author: desiredSubmision.GetAuthor(),
		Title:  desiredSubmision.GetTitle(),
		Link:   desiredSubmision.Data.Url,
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
