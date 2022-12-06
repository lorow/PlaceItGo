package utils

import (
	"errors"
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
