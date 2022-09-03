package internal

import (
	"context"
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

	ctx := r.Context()

	if width_err != nil || height_err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageService, ok := ctx.Value("imageService").(*ImageManager)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	imageData, err := imageService.GetImage(animal, width, height)

	if err != nil {
		log.Error().Msgf("Something went wrong while retrieving the image: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(imageData.data)
}

func imageServiceContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config, config_err := GetConfig()
		if config_err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		redisCache, redisConnectionError := GetRedisCache(*config)

		if redisConnectionError != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		imageService := ImageManager{
			redisCache: redisCache,
		}

		ctx := context.WithValue(r.Context(), "imageService", imageService)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func StartServer() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", serveIndex)
	router.Route("/{animal}/{width}/{height}", func(r chi.Router) {
		r.Use(imageServiceContext)
		r.Get("/", getImage)
	})

	port := ":8080"

	log.Printf("Started serving on port http://localhost%s", port)
	return http.ListenAndServe(port, router)
}
