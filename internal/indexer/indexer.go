package indexer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/Sasikuttan2163/Telescope/internal/db/qdrantdb"
	"github.com/Sasikuttan2163/Telescope/internal/embed"
	"github.com/Sasikuttan2163/Telescope/internal/transport"
)

func IndexStar(ctx context.Context, ollamaConfig config.OllamaConfig, starConfig config.StarConfig, qdrant *qdrantdb.Qdrant, collectionName string) error {

	fetchCtx, cancel := context.WithTimeout(ctx, time.Duration(starConfig.Timeout)*time.Second)
	defer cancel()

	tools, err := transport.FetchToolsOfStar(fetchCtx, starConfig)
	if err != nil {
		return fmt.Errorf("Failed to fetch tools for %s: %w", starConfig.ID, err)
	}

	embedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = embed.OllamaGetToolVector(embedCtx, ollamaConfig.Host, ollamaConfig.Port, &ollamaConfig.NumGpu, ollamaConfig.Model, starConfig.Name, &tools)
	if err != nil {
		return fmt.Errorf("Failed to get ollama vectors for MCP %s: %w", starConfig.ID, err)
	}

	insertCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = qdrant.BatchInsert(insertCtx, collectionName, tools)
	return err
}

func IndexAllStars(ctx context.Context, mainConfig config.MainConfig) (successCount int, errs []error) {
	indexStarsCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		indexWg sync.WaitGroup
	)

	qdrantClient, err := qdrantdb.NewQdrant(mainConfig.Qdrant.Host, mainConfig.Qdrant.Port)
	qdrantClient.CreateCollection(indexStarsCtx, mainConfig.Qdrant.CollectionName, mainConfig.Ollama.EmbedDim)

	if err != nil {
		errs = append(errs, err)
		return
	}

	for _, star := range mainConfig.Stars {
		indexWg.Add(1)

		go func(s config.StarConfig) {
			defer indexWg.Done()

			cctx, cancel := context.WithTimeout(
				indexStarsCtx,
				time.Duration(star.Timeout)*time.Second,
			)
			defer cancel()

			indexError := IndexStar(cctx, mainConfig.Ollama, star, qdrantClient, mainConfig.Qdrant.CollectionName)
			if indexError != nil {
				errs = append(errs, indexError)
			} else {
				successCount += 1
			}
		}(star)
	}

	indexWg.Wait()
	return
}
