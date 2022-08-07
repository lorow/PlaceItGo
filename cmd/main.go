package main

import (
	"github.com/lorow/placeitgo/api"
	log "github.com/sirupsen/logrus"
)

func main() {
	mainLogger := log.New()
	logger := log.NewEntry(mainLogger)

	err := api.StartServer()

	if err != nil {
		 panic("couldn't start the API server up")
	}
	logger.Info("Started serving")
}