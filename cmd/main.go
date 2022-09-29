package main

import (
	"fmt"

	pkg "github.com/lorow/placeitgo/internal"
	"github.com/rs/zerolog"
)

func main() {
	config, config_err := pkg.GetConfig()
	if config_err != nil {
		fmt.Println("Couldn't load config")
	}

	redisCache, redisConnectionError := pkg.NewRedisCache(config)
	if redisConnectionError != nil {
		fmt.Printf("Couldn't connect to redis instance - {%s} \n", redisConnectionError)
	} else {
		fmt.Println("Successfully connected to redis!")
	}

	redditService, err := pkg.NewRedditService(config)
	if err != nil {
		panic(fmt.Sprintf("Could not start the reddit service %s", err))
	}

	imageService := pkg.ImageManager{
		RedisCache:    redisCache,
		RedditService: redditService,
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	server_err := pkg.StartServer(imageService)

	if server_err != nil {
		panic("couldn't start the API server up!")
	}
}
