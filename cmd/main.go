package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/Sasikuttan2163/Telescope/internal/indexer"
)

func main() {
	config, err := config.GetConfig("test.json")
	if err != nil {
		log.Fatal("Fatal error " + err.Error())
	}

	succ, errs := indexer.IndexAllStars(context.Background(), config)
	fmt.Printf("Successfully indexed: %d out of %d\n", succ, len(config.Stars))

	for _, err := range errs {
		fmt.Printf("%s", err.Error())
	}
}
