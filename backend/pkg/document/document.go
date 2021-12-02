package document

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"gritter/pkg/context"
)

var (
	ErrNotFound = errors.New("document not found")
)

type Document interface {
	Search(ctx context.Context, name Name, offset, limit int64, query bson.M, sort bson.D, doc interface{}) error
	CreateOne(ctx context.Context, name Name, doc interface{}) error
	GetOne(ctx context.Context, name Name, query bson.M, doc interface{}) error
	UpdateOne(ctx context.Context, name Name, query bson.M, doc interface{}) error
	DeleteOne(ctx context.Context, name Name, query bson.M) error
}
