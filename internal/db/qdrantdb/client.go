package qdrantdb

import (
	"context"
	"reflect"

	"github.com/Sasikuttan2163/Telescope/internal/types"
	"github.com/qdrant/go-client/qdrant"
)

type Qdrant struct {
	Client *qdrant.Client
}

func NewQdrant(addr string, port int) (*Qdrant, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: addr,
		Port: port,
	})

	if err != nil {
		return nil, err
	}

	return &Qdrant{Client: client}, nil
}

func (q *Qdrant) CreateCollection(ctx context.Context, collectionName string, embedDim uint64) error {
	return q.Client.CreateCollection(
		ctx,
		&qdrant.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     embedDim,
				Distance: qdrant.Distance_Cosine,
			}),
		},
	)
}

func (q *Qdrant) Insert(ctx context.Context, collectionName string, id string, vectors []float32, payload map[string]any) error {
	_, err := q.Client.Upsert(
		ctx,
		&qdrant.UpsertPoints{
			CollectionName: collectionName,
			Points: []*qdrant.PointStruct{
				{
					Id:      qdrant.NewID(id),
					Vectors: qdrant.NewVectors(vectors...),
					Payload: qdrant.NewValueMap(payload),
				},
			},
		},
	)
	return err
}

func structToMap(s any, excludeFields ...string) map[string]any {
	exclude := make(map[string]bool)
	for _, f := range excludeFields {
		exclude[f] = true
	}

	v := reflect.ValueOf(s).Elem()
	result := make(map[string]any)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if exclude[field.Name] {
			continue
		}
		result[field.Name] = v.Field(i).Interface()
	}

	return result
}

func (q *Qdrant) BatchInsert(ctx context.Context, collectionName string, tools []*types.Tool) error {
	points := make([]*qdrant.PointStruct, len(tools))

	for i, v := range tools {
		payload := structToMap(v, "Vector")
		points[i] = &qdrant.PointStruct{
			Id:      qdrant.NewID(v.Uuid),
			Vectors: qdrant.NewVectors(v.Vector...),
			Payload: qdrant.NewValueMap(payload),
		}
	}

	_, err := q.Client.Upsert(
		ctx,
		&qdrant.UpsertPoints{
			CollectionName: collectionName,
			Points:         points,
		},
	)
	return err
}

func (q *Qdrant) Query(ctx context.Context, collectionName string, query []float32) ([]*qdrant.ScoredPoint, error) {
	return q.Client.Query(
		ctx,
		&qdrant.QueryPoints{
			CollectionName: collectionName,
			Query:          qdrant.NewQuery(query...),
			WithPayload:    qdrant.NewWithPayload(true),
		},
	)
}

func (q *Qdrant) GetAllPoints(ctx context.Context, collectionName string) ([]*qdrant.RetrievedPoint, error) {
	var allPoints []*qdrant.RetrievedPoint
	var nextOffset *qdrant.PointId

	for {
		resp, newOffset, err := q.Client.ScrollAndOffset(ctx, &qdrant.ScrollPoints{
			CollectionName: collectionName,
			WithPayload:    qdrant.NewWithPayload(true),
			WithVectors:    qdrant.NewWithVectors(false),
			Offset:         nextOffset,
		})
		if err != nil {
			return nil, err
		}

		allPoints = append(allPoints, resp...)

		if len(resp) == 0 || newOffset == nil {
			break
		}

		nextOffset = newOffset
	}

	return allPoints, nil
}
