package internal

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Here will be a template"))
}

func getImage(w http.ResponseWriter, r *http.Request) {

	animal := chi.URLParam(r, "animal")
	width_str := chi.URLParam(r, "width")
	height_str := chi.URLParam(r, "height")

	width, width_err := strconv.Atoi(width_str)
	height, height_err := strconv.Atoi(height_str)

	if width_err != nil || height_err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// this should be in context
	config, config_err := GetConfig()
	if config_err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// this should be in context
	redisCache := GetRedisCache(*config)
	imageService := ImageManager{
		redisCache,
	}

	imageData, err := imageService.GetImage(animal, width, height)

	if err != nil {
		log.Error().Msgf("Something went wrong while retrieving the image: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(imageData)
}

func StartServer() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", serveIndex)
	router.Get("/{animal}/{width}/{height}", getImage)

	port := ":8080"

	log.Printf("Started serving on port http://localhost%s", port)
	return http.ListenAndServe(port, router)
}
