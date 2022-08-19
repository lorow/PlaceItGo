package internal

import (
	"fmt"
)

type ImageService interface {
	GetImage(animal string, width, height int) ([]byte, error)
}

type ImageManager struct {
	redisCache *RedisCache
}

func (i ImageManager) GetImage(animal string, width, height int) ([]byte, error) {
	// todo, try and see if it was reddis connection issue (maybe add a preflight check while getting redis?)
	// todo, add getting images from reddit
	// todo, add randomization to retrieved images
	cahedImages, err := i.redisCache.GetImages(fmt.Sprintf("%s_%xx%x", animal, width, height))
	return cahedImages[0], err
}
