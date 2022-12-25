package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func ConvertResolutionFormString(resolution, delimiter string) (int, int, error) {
	dimensionsSplit := strings.Split(resolution, delimiter)

	if len(dimensionsSplit) == 1 {
		return -1, -1, errors.New("invalid resolution format")
	}

	width, err := strconv.Atoi(dimensionsSplit[0])
	if err != nil {
		return -1, -1, errors.New("invalid resolution value")
	}

	height, err := strconv.Atoi(dimensionsSplit[1])
	if err != nil {
		return -1, -1, errors.New("invalid resolution value")
	}

	return width, height, nil
}

func FetchImageFromURL(link string) ([]byte, error) {
	response, err := http.Get(link)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("error getting image data: %s", err)
	}

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return imageData, nil
}
