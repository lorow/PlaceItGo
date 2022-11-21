package internal

import (
	"strconv"
	"strings"
)

func ConvertResolutionFormString(resolution, delimiter string) (int, int, error) {
	dimensionsSplit := strings.Split(resolution, delimiter)

	width, err := strconv.Atoi(dimensionsSplit[0])
	if err != nil {
		return -1, -1, err
	}

	height, err := strconv.Atoi(dimensionsSplit[1])
	if err != nil {
		return -1, -1, err
	}

	return width, height, nil
}
