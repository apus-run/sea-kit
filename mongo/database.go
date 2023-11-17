package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is a handle to a MongoDB database
type Database struct {
	database *mongo.Database

	registry *bsoncodec.Registry
}

func (d *Database) Client() *Client {
	return &Client{
		client:   d.database.Client(),
		registry: d.registry,
	}
}
func (d *Database) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *SingleResult {
	option := options.RunCmd()
	if len(opts) > 0 && opts != nil {
		option = opts[0]
	}
	return d.database.RunCommand(ctx, runCommand, option)
}
func (d *Database) RunCommandCursor(ctx context.Context, cmd interface{}) *Cursor {
	cur, err := d.database.RunCommandCursor(ctx, cmd)
	return &Cursor{
		ctx:    ctx,
		cursor: cur,
		err:    err,
	}
}

// Collection gets collection from database
func (d *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	coll := d.database.Collection(name, opts...)

	return &Collection{
		collection: coll,
		registry:   d.registry,
	}
}

func (d *Database) CreateCollection(ctx context.Context, name string, opts *options.CollectionOptions) (*Collection, error) {
	if err := d.database.CreateCollection(ctx, name); err != nil {
		return nil, err
	}
	cp := d.database.Collection(name, opts)
	return &Collection{
		collection: cp,
		registry:   d.registry,
	}, nil
}

// GetDatabaseName returns the name of database
func (d *Database) GetDatabaseName() string {
	return d.database.Name()
}

// DropDatabase drops database
func (d *Database) DropDatabase(ctx context.Context) error {
	return d.database.Drop(ctx)
}
