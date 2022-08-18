package main

import (
	"io"
	"log"
	"os"

	pkg "github.com/lorow/placeitgo/internal"
)

func setupLogging() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
			panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func main() {
	setupLogging()

	server_err := pkg.StartServer()

	if server_err != nil {
		panic("couldn't start the API server up")
	}

}
