package imageservice

import (
	"math/rand"
	"placeitgo/model"
	"placeitgo/storage"
	"placeitgo/utils"
	"time"
)

type ImageDownloader interface {
	GetImages(animal string, width, height int) ([]model.ImageDBEntry, error)
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
	var downloadedImageData []byte
	// todos
	// add cropping
	// add tests
	imageEntry, err := i.storage.GetImage(animal, width, height)
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

	imageEntries, err := i.downloader.GetImages(animal, width, height)
	if err != nil {
		return model.ImageResponse{}, err
	}

	// before anything, store the image entries
	i.storage.SaveImageEntries(imageEntries, animal)

	// then select a random entry and work with it
	rand.Seed(time.Now().Unix())
	imageEntry = imageEntries[rand.Intn(len(imageEntries))]
	downloadedImageData, err = i.fetchImage(imageEntry)
	if err != nil {
		return model.ImageResponse{}, err
	}

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
