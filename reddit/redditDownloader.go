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

// GetImages Downloads all the matching images from reddit,
// whose dimensions are bigger or matching the ones provided by the user.
// Only up to three most similar in terms of resolution
func (r RedditService) GetImages(animal string, width, height int) ([]model.ImageDBEntry, error) {
	var maxRetries = 3
	var matchingSubmissions []model.ImageDBEntry
	var lastSubmission *mira.PostListingChild = nil

	submissions, err := r.bot.GetSubredditPosts(animal, "hot", "all", 10)
	if err != nil {
		return []model.ImageDBEntry{}, err
	}

	for i := 0; i < maxRetries; i++ {
		if lastSubmission != nil {
			submissions, err = r.bot.GetSubredditPostsAfter(animal, lastSubmission.After, 10)
			if err != nil {
				return []model.ImageDBEntry{}, err
			}
		}
		r.processSubmissions(&submissions, &matchingSubmissions, width, height)

		if len(matchingSubmissions) < 3 {
			lastSubmission = &submissions[len(submissions)-1]
		} else {
			lastSubmission = nil
		}
	}

	if matchingSubmissions == nil {
		return []model.ImageDBEntry{}, fmt.Errorf("could not find desired placeholder image for %s %dx%d", animal, width, height)
	}

	return r.filterMatchingSubmissions(&matchingSubmissions), nil
}

func (r RedditService) processSubmissions(submissions *[]mira.PostListingChild, matchingSubmissions *[]model.ImageDBEntry, width, height int) {
	for _, submission := range *submissions {
		if !submission.Data.Preview.Enabled {
			continue
		}
		submissionWidth := submission.Data.Preview.Images[0].Source.Width
		submissionHeight := submission.Data.Preview.Images[0].Source.Height

		matchingWidth := bool(int(submissionWidth)-width >= 0)
		matchingHeight := bool(int(submissionHeight)-height >= 0)

		if matchingWidth && matchingHeight {
			matchingSubmission := model.ImageDBEntry{
				Author: submission.GetAuthor(),
				Title:  submission.GetTitle(),
				Link:   submission.Data.Url,
				Width:  int(submissionWidth),
				Height: int(submissionHeight),
			}

			*matchingSubmissions = append(*matchingSubmissions, matchingSubmission)
		}
	}
}

func (r RedditService) filterMatchingSubmissions(matchingSubmissions *[]model.ImageDBEntry) []model.ImageDBEntry {
	return []model.ImageDBEntry{}
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
