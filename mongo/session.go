package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session = mongo.Session
type SessionContext = mongo.SessionContext

type session struct {
	mongo.Session
}

var _ mongo.Session = (*session)(nil)

func (ws *session) EndSession(ctx context.Context) {
	ws.Session.EndSession(ctx)
}

func (ws *session) StartTransaction(topts ...*options.TransactionOptions) error {
	return ws.Session.StartTransaction(topts...)
}

func (ws *session) AbortTransaction(ctx context.Context) error {
	return ws.Session.AbortTransaction(ctx)
}

func (ws *session) CommitTransaction(ctx context.Context) error {
	return ws.Session.CommitTransaction(ctx)
}

func (ws *session) ClusterTime() (raw bson.Raw) {
	return ws.Session.ClusterTime()
}

func (ws *session) AdvanceClusterTime(br bson.Raw) error {
	return ws.Session.AdvanceClusterTime(br)
}

func (ws *session) OperationTime() (ts *primitive.Timestamp) {
	return ws.Session.OperationTime()
}

func (ws *session) AdvanceOperationTime(pt *primitive.Timestamp) error {
	return ws.Session.AdvanceOperationTime(pt)
}

func (ws *session) WithTransaction(ctx context.Context, fn func(sessCtx SessionContext) (interface{}, error),
	opts ...*options.TransactionOptions) (out interface{}, err error) {
	return ws.Session.WithTransaction(ctx, fn, opts...)
}
