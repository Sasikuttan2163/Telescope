package main

import (
	"log"

	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/Sasikuttan2163/Telescope/internal/transport"
)

func main() {
	config, err := config.GetConfig("test.json")
	if err != nil {
		log.Fatal("Fatal error " + err.Error())
	}
	transport.FetchTools(config)
}
