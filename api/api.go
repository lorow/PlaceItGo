package api

import (
	"net/http"
	"placeitgo/model"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ImageService interface {
	GetImage(animal string, width, height int) (model.ImageResponse, error)
}

type PlaceItGoHandler struct {
	imageService ImageService
}

func (p PlaceItGoHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Here will be a template"))
}

func (p PlaceItGoHandler) getImage(w http.ResponseWriter, r *http.Request) {

	animal := chi.URLParam(r, "animal")
	widthStr := chi.URLParam(r, "width")
	heightStr := chi.URLParam(r, "height")

	width, widthErr := strconv.Atoi(widthStr)
	height, heightErr := strconv.Atoi(heightStr)

	if widthErr != nil || heightErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageData, err := p.imageService.GetImage(animal, width, height)

	if err != nil {
		log.Error().Msgf("Something went wrong while retrieving the image: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(imageData.Data)
}

func StartServer(imageManager ImageService) error {

	handler := PlaceItGoHandler{
		imageService: imageManager,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", handler.serveIndex)
	router.Route("/{animal}/{width}/{height}", func(r chi.Router) {
		r.Get("/", handler.getImage)
	})

	port := ":8080"

	log.Printf("Started serving on port %s", port)
	return http.ListenAndServe(port, router)
}
