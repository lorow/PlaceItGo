package internal

import (
	"fmt"
	"math/rand"
)

type ImageService interface {
	GetImage(animal string, width, height int) (RedditImage, error)
}

type ImageManager struct {
	RedisCache    *RedisCache
	RedditService *RedditService
}

type RedditImage struct {
	name string
	data []byte
}

func (i ImageManager) GetImage(animal string, width, height int) (RedditImage, error) {
	// todo, add getting images from reddit
	images, err := i.RedisCache.GetImages(fmt.Sprintf("%s_%xx%x", animal, width, height))
	if err == nil && len(images) != 0 {
		randomIndex := rand.Intn(len(images))
		return images[randomIndex], nil
	}

	image, err := i.RedditService.GetImage(animal, width, height)
	if err == nil {
		i.RedisCache.SaveImage(image.name, image.data)
		return image, nil
	}

	return RedditImage{}, err
}
