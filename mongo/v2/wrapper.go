package mongo

import (
	"context"
	"fmt"
)

// Mongo is the MongoDB
type Mongo struct {
	*Client
	*Database
	*Collection
}

func Open(ctx context.Context, opts ...Option) (*Mongo, error) {
	cfg := Apply(opts...)
	client, err := client(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new client fail: %v", err)
	}

	db := client.Database(cfg.DatabaseName)
	coll := db.Collection(cfg.CollectionName)

	return &Mongo{
		Client:     client,
		Database:   db,
		Collection: coll,
	}, nil
}
