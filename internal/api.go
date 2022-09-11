package internal

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type PlaceItGoHandler struct {
	imageService ImageService
}

func (p PlaceItGoHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Here will be a template"))
}

func (p PlaceItGoHandler) getImage(w http.ResponseWriter, r *http.Request) {

	animal := chi.URLParam(r, "animal")
	width_str := chi.URLParam(r, "width")
	height_str := chi.URLParam(r, "height")

	width, width_err := strconv.Atoi(width_str)
	height, height_err := strconv.Atoi(height_str)

	if width_err != nil || height_err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageData, err := p.imageService.GetImage(animal, width, height)

	if err != nil {
		log.Error().Msgf("Something went wrong while retrieving the image: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(imageData.data)
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
