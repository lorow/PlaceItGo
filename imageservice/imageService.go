package imageservice

import (
	"placeitgo/model"
	"placeitgo/storage"
	"placeitgo/utils"
)

type ImageDownloader interface {
	GetImage(animal string, width, height int) (model.ImageDBEntry, error)
}

type ImageService interface {
	GetImage(animal string, width, height int) (model.ImageResponse, error)
}

type ImageProcessor interface {
	ProcessImageEntry(author, title string, data []byte) ([]byte, error)
}

type ImageHandler struct {
	storage    storage.Storage
	downloader ImageDownloader
	processor  ImageProcessor
}

func (i ImageHandler) GetImage(animal string, width, height int) (model.ImageResponse, error) {
	imageEntry, err := i.storage.GetImage(width, height, animal)
	if err == nil {
		downloadedImageData, err := i.fetchImage(imageEntry)
		if err != nil {
			return model.ImageResponse{}, err
		}

		processedImage, err := i.processor.ProcessImageEntry(imageEntry.Author, imageEntry.Title, downloadedImageData)
		if err != nil {
			return model.ImageResponse{}, err
		}

		return model.ImageResponse{
			Data: processedImage,
		}, nil
	}

	imageEntry, err = i.downloader.GetImage(animal, width, height)
	if err != nil {
		return model.ImageResponse{}, err
	}

	downloadedImageData, err := i.fetchImage(imageEntry)
	if err != nil {
		return model.ImageResponse{}, err
	}

	// before anything, store the fresh image
	i.storage.SaveImage(imageEntry.Height, imageEntry.Width, imageEntry.Author, imageEntry.Title, animal, imageEntry.Link)

	processedImage, err := i.processor.ProcessImageEntry(imageEntry.Author, imageEntry.Title, downloadedImageData)
	if err != nil {
		return model.ImageResponse{}, nil
	}

	return model.ImageResponse{
		Data: processedImage,
	}, nil
}

func (i ImageHandler) fetchImage(entry model.ImageDBEntry) ([]byte, error) {
	return utils.FetchImageFromURL(entry.Link)
}

func NewImageService(storage storage.Storage, downloader ImageDownloader, processor ImageProcessor) ImageHandler {
	return ImageHandler{
		storage:    storage,
		downloader: downloader,
		processor:  processor,
	}
}
