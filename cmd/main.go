package main

import (
	pkg "github.com/lorow/placeitgo/internal"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	server_err := pkg.StartServer()

	if server_err != nil {
		panic("couldn't start the API server up!")
	}

}
