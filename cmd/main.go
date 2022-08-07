package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	mainLogger := log.New()
	logger := log.NewEntry(mainLogger)
}