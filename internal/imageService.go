package internal

import (
	"fmt"
	"math/rand"
)

type ImageService interface {
	GetImage(animal string, width, height int) ([]byte, error)
}

type ImageManager struct {
	redisCache    *RedisCache
	redditService *RedditService
}

type RedditImage struct {
	name string
	data []byte
}

func (i ImageManager) GetImage(animal string, width, height int) (RedditImage, error) {
	// todo, add getting images from reddit

	images, err := i.redisCache.GetImages(fmt.Sprintf("%s_%xx%x", animal, width, height))
	if err == nil || len(images) == 0 {
		randomIndex := rand.Intn(len(images))
		return images[randomIndex], nil
	}

	image, err := i.redditService.GetImage(animal, width, height)
	if err == nil {
		i.redisCache.SaveImage(image.name, image.data)
		return image, nil
	}

	return RedditImage{}, err
}
