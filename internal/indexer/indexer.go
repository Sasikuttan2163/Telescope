package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/Sasikuttan2163/Telescope/internal/db/qdrantdb"
	"github.com/Sasikuttan2163/Telescope/internal/embed"
	"github.com/Sasikuttan2163/Telescope/internal/transport"
	"github.com/Sasikuttan2163/Telescope/internal/types"
	"github.com/qdrant/go-client/qdrant"
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

func GetAllIndexedStars(ctx context.Context, mainConfig config.MainConfig) ([]types.Tool, error) {
	qdrantClient, err := qdrantdb.NewQdrant(mainConfig.Qdrant.Host, mainConfig.Qdrant.Port)
	if err != nil {
		return nil, err
	}

	points, err := qdrantClient.GetAllPoints(ctx, mainConfig.Qdrant.CollectionName)
	if err != nil {
		return nil, err
	}

	return payloadsToTools(points)
}

func GetTopKTools(ctx context.Context, mainConfig config.MainConfig, queryVector []float32) ([]types.Tool, error) {
	qdrantClient, err := qdrantdb.NewQdrant(mainConfig.Qdrant.Host, mainConfig.Qdrant.Port)
	if err != nil {
		return nil, err
	}

	points, err := qdrantClient.Query(ctx, mainConfig.Qdrant.CollectionName, queryVector)
	if err != nil {
		return nil, err
	}

	return payloadsToTools(points)
}

func payloadsToTools[T any](points []T) ([]types.Tool, error) {
	tools := make([]types.Tool, 0, len(points))

	for _, point := range points {
		payloadMap := getPayload(point)
		if payloadMap == nil {
			continue
		}

		payloadBytes, err := json.Marshal(payloadMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}

		var tool types.Tool
		if err := json.Unmarshal(payloadBytes, &tool); err != nil {
			return nil, fmt.Errorf("failed to unmarshal into tool: %w", err)
		}

		tools = append(tools, tool)
	}

	return tools, nil
}

func getPayload[T any](point T) map[string]any {
	v := any(point)
	switch p := v.(type) {
	case *qdrant.RetrievedPoint:
		return qdrantPayloadToMap(p.Payload)
	case *qdrant.ScoredPoint:
		return qdrantPayloadToMap(p.Payload)
	default:
		return nil
	}
}

func qdrantPayloadToMap(payload map[string]*qdrant.Value) map[string]any {
	result := make(map[string]any)
	for k, v := range payload {
		result[k] = qdrantValueToAny(v)
	}
	return result
}

func qdrantValueToAny(v *qdrant.Value) any {
	switch kind := v.GetKind().(type) {
	case *qdrant.Value_StringValue:
		return kind.StringValue
	case *qdrant.Value_IntegerValue:
		return kind.IntegerValue
	case *qdrant.Value_DoubleValue:
		return kind.DoubleValue
	case *qdrant.Value_BoolValue:
		return kind.BoolValue
	case *qdrant.Value_ListValue:
		list := make([]any, 0, len(kind.ListValue.Values))
		for _, item := range kind.ListValue.Values {
			list = append(list, qdrantValueToAny(item))
		}
		return list
	case *qdrant.Value_StructValue:
		return qdrantPayloadToMap(kind.StructValue.Fields)
	default:
		return nil
	}
}
