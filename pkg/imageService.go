package pkg

type ImageService interface {
	GetImage(animal string, width, height int) ([]byte, error)
}

type ImageManager struct{}

func (i ImageManager) GetImage(animal string, width, height int) ([]byte, error) {
	return nil, nil
}
