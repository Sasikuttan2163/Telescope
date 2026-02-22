// Package qdrantdb provides a Qdrant vector database client.
package qdrantdb

import (
	"context"

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

func (q *Qdrant) Query(ctx context.Context, collectionName string, query []float32) ([]*qdrant.ScoredPoint, error) {
	return q.Client.Query(
		ctx,
		&qdrant.QueryPoints{
			CollectionName: collectionName,
			Query:          qdrant.NewQuery(query...),
		},
	)
}
