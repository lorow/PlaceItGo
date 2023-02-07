package main

import (
	"fmt"
	"placeitgo/api"
	"placeitgo/config"
	"placeitgo/imageprocessor"
	"placeitgo/imageservice"
	"placeitgo/reddit"
	"placeitgo/storage"

	"github.com/rs/zerolog"
)

func main() {
	config, configErr := config.GetConfig()
	if configErr != nil {
		fmt.Println("Couldn't load config")
	}

	redisCache, redisConnectionError := storage.NewRedisCache(config)
	if redisConnectionError != nil {
		panic(fmt.Sprintf("Couldn't connect to redis instance - {%s} \n", redisConnectionError))
	} else {
		fmt.Println("Successfully connected to redis!")
	}
	// TODO refactor this later to a more generic service
	redditDownloader, err := reddit.NewRedditDownloader(config)
	if err != nil {
		panic(fmt.Sprintf("Could not start the reddit service %s", err))
	}

	imageProcessor := imageprocessor.GetNewImageProcessor()

	imageService := imageservice.NewImageService(redisCache, redditDownloader, imageProcessor)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	serverErr := api.StartServer(imageService)

	if serverErr != nil {
		panic("couldn't start the API server up!")
	}
}
