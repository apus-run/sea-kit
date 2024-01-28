package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client creates client to mongo
type Client struct {
	client *mongo.Client

	registry *bsoncodec.Registry
}

// NewClient creates MongoDB client
func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := Apply(opts...)
	return client(ctx, cfg)
}

func client(ctx context.Context, cfg *Config) (*Client, error) {
	option := options.Client()
	option = cfg.ClientOptions
	option.ApplyURI(cfg.Uri)

	client, err := mongo.Connect(ctx, option)
	if err != nil {
		return nil, fmt.Errorf("can't connect to mongo: %w", err)
	}

	// half of default connect timeout
	pCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err = client.Ping(pCtx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("can't ping mongo: %w", err)
	}

	return &Client{
		client:   client,
		registry: option.Registry,
	}, nil
}

// Disconnect closes sockets to the topology referenced by this Client.
func (c *Client) Disconnect(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	return err
}

// Close closes sockets to the topology referenced by this Client.
func (c *Client) Close(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	return err
}

// Ping confirm connection is alive
func (c *Client) Ping(timeout int64) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if err = c.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}

// Database create connection to database
func (c *Client) Database(name string, dbOpts ...*options.DatabaseOptions) *Database {
	db := c.client.Database(name, dbOpts...)
	return &Database{
		database: db,
		registry: c.registry,
	}
}

// ServerVersion get the version of mongoDB server, like 4.4.0
func (c *Client) ServerVersion() string {
	var buildInfo bson.Raw
	err := c.client.Database("admin").RunCommand(
		context.Background(),
		bson.D{{"buildInfo", 1}},
	).Decode(&buildInfo)
	if err != nil {
		fmt.Println("run command err", err)
		return ""
	}
	v, err := buildInfo.LookupErr("version")
	if err != nil {
		fmt.Println("look up err", err)
		return ""
	}
	return v.StringValue()
}
