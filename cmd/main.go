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

	redisCache, redisConnectionError := pkg.NewRedisCache(*config)
	if redisConnectionError != nil {
		fmt.Println("Couldn't connect to redis instance")
	}

	redditService := pkg.NewRedditService()
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
