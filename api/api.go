package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func serveIndex(writer http.ResponseWriter, request *http.Request){
	writer.Write([]byte("welcome, here will be a basic template"))
}


func getImage(writer http.ResponseWriter, request *http.Request) {

	animal  := chi.URLParam(request, "animal")
	width  := chi.URLParam(request, "width")
	height  := chi.URLParam(request, "height")

	writer.Write([]byte(fmt.Sprintf("%s %s %s", animal, width, height)))
}

func StartServer() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", serveIndex)
	router.Get("/{animal}/{width}/{height}", getImage)

	err := http.ListenAndServe(":8080", router)

	return err
}