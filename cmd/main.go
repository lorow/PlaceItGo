package main

import (
	"fmt"
	"placeitgo/api"
	"placeitgo/config"
	"placeitgo/reddit"
	"placeitgo/storage"

	"github.com/rs/zerolog"
)

func main() {
	config, config_err := config.GetConfig()
	if config_err != nil {
		fmt.Println("Couldn't load config")
	}

	redisCache, redisConnectionError := storage.NewRedisCache(config)
	if redisConnectionError != nil {
		panic(fmt.Sprintf("Couldn't connect to redis instance - {%s} \n", redisConnectionError))
	} else {
		fmt.Println("Successfully connected to redis!")
	}

	imageService, err := reddit.NewRedditService(config, redisCache)
	if err != nil {
		panic(fmt.Sprintf("Could not start the reddit service %s", err))
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	server_err := api.StartServer(imageService)

	if server_err != nil {
		panic("couldn't start the API server up!")
	}
}
