package reddit

import (
	"fmt"
	"placeitgo/config"
	"placeitgo/storage"

	"github.com/thecsw/mira"
)

type RedditService struct {
	bot   *mira.Reddit
	cache *storage.RedisCache
}

func (r RedditService) GetImage(animal string, width, height int) ([]byte, error) {
	// todo, add getting images from reddit

	// here's how it shoudl function.
	// request:
	//   -> check for this exact dimensions and animal in cache:
	//      -> found:
	//        -> fetch image info
	//        -> fetch info about author
	//        -> fetch info about title
	//        -> process the image
	//        -> return the image
	//      -> not found:
	//        -> check for higher dimension:
	//           -> get all keys that include animal name
	//           -> select random key
	//           -> get data from under that key:
	//              -> fetch info about author
	//              -> fetch info about title
	//              -> process the image
	//              -> return the image
	//         -> nothing fits
	//			  -> get reddit results
	//			  -> for every result:
	//			     -> make note of the title and author
	//			     -> download image
	//			     -> check dimensions
	//			     -> dimensions are a close match:
	//				 -> save the result
	//               -> process the image
	//               -> return the image

	// NOTE: Instead of saving the image, save the link to it and process it on demand

	imageEntry, err := r.cache.GetImage(width, height, animal)

	if err == nil {
		// we got the iamge from cache, let's download it, process it, and return
		freshImage, err := r.fetchImage(imageEntry)
		if err != nil {
			return []byte{},
				err
		}
		processedImage := r.processImageEntry(imageEntry.Author, imageEntry.Title, freshImage)
		return processedImage, nil
	} else {
		// something went wrong with getting the image from cache, possibly nothing found
		// get it straight from reddit
		return []byte{}, nil
	}
}

func (r RedditService) fetchImage(entry storage.ImageEntry) ([]byte, error) {
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
