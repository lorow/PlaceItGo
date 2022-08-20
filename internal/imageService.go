package internal

import (
	"fmt"
	"math/rand"
)

type ImageService interface {
	GetImage(animal string, width, height int) ([]byte, error)
}

type ImageManager struct {
	redisCache *RedisCache
}

func (i ImageManager) GetImage(animal string, width, height int) ([]byte, error) {
	// todo, add getting images from reddit
	// todo, add randomization to retrieved images
	images, err := i.redisCache.GetImages(fmt.Sprintf("%s_%xx%x", animal, width, height))
	randomIndex := rand.Intn(len(images))
	image := images[randomIndex]

	return image, err
}
