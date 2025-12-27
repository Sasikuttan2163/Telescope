package main

import (
	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/Sasikuttan2163/Telescope/internal/transport"
)

func main() {
	transport.FetchTools(config.GetConfig("test.json"))
}
