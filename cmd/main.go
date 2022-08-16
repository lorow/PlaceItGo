package main

import (
	pkg "github.com/lorow/placeitgo/internal"
	log "github.com/sirupsen/logrus"
)

func main() {
	mainLogger := log.New()
	logger := log.NewEntry(mainLogger)

	err := pkg.StartServer()

	if err != nil {
		panic("couldn't start the API server up")
	}
	logger.Info("Started serving")
}
