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

func (r RedditService) GetImages(animal string, width, height int) ([]model.ImageDBEntry, error) {
	var matchingSubmissions []model.ImageDBEntry

	// todo add max tries
	submissions, err := r.bot.GetSubredditPosts(animal, "hot", "all", 10)

	if err != nil {
		return []model.ImageDBEntry{}, err
	}

	for _, submission := range submissions {
		if !submission.Data.Preview.Enabled {
			continue
		}
		submissionWidth := submission.Data.Preview.Images[0].Source.Width
		submissionHeight := submission.Data.Preview.Images[0].Source.Height

		matchingWidth := bool(int(submissionWidth)-width > 0 && float32(int(submissionWidth)-width) <= float32(width)*0.3)
		matchingHeight := bool(int(submissionHeight)-height > 0 && float32(int(submissionHeight)-height) <= float32(height)*0.3)

		if matchingWidth && matchingHeight {
			matchingSubmission := model.ImageDBEntry{
				Author: submission.GetAuthor(),
				Title:  submission.GetTitle(),
				Link:   submission.Data.Url,
				Width:  int(submissionWidth),
				Height: int(submissionHeight),
			}

			matchingSubmissions = append(matchingSubmissions, matchingSubmission)
		}
	}

	if matchingSubmissions == nil {
		return []model.ImageDBEntry{}, fmt.Errorf("could not find desired placeholder iamge for %s %dx%d", animal, width, height)
	}

	return matchingSubmissions, nil
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
