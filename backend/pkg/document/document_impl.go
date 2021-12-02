package document

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"gritter/pkg/context"
	"gritter/pkg/db"
)

var (
	dbName = "gritter"
)

func NewDocument() Document {
	return &impl{
		mongo: db.MustGetMongo(),
	}
}

type impl struct {
	mongo *mongo.Client
}

func (im *impl) Search(ctx context.Context, name Name, offset, limit int64, query bson.M, sort bson.D, doc interface{}) error {
	option := options.Find().SetSkip(offset).SetLimit(limit)
	if sort != nil {
		option.SetSort(sort)
	}
	cursor, err := im.mongo.Database(dbName).Collection(name.String()).Find(ctx, query, option)
	if err != nil {
		ctx.With(zap.Error(err)).Error("mongo.Collection.Find failed in document.Document.Search")
		return err
	}

	if err := cursor.All(ctx, doc); err != nil {
		ctx.With(zap.Error(err)).Error("mongo.Cursor.All failed in document.Document.Search")
		return err
	}

	return nil
}

func (im *impl) CreateOne(ctx context.Context, name Name, doc interface{}) error {
	_, err := im.mongo.Database(dbName).Collection(name.String()).InsertOne(ctx, doc)
	if err != nil {
		ctx.With(zap.Error(err)).Error("mongo.Collection.InsertOne failed in document.Document.CreateOne")
		return err
	}

	return nil
}

func (im *impl) GetOne(ctx context.Context, name Name, query bson.M, doc interface{}) error {
	res := im.mongo.Database(dbName).Collection(name.String()).FindOne(ctx, query)
	err := res.Err()
	if err == mongo.ErrNoDocuments {
		return ErrNotFound
	} else if err != nil {
		ctx.With(zap.Error(err)).Error("mongo.Collection.FindOne failed in document.Document.GetOne")
		return err
	}

	if err := res.Decode(doc); err != nil {
		ctx.With(zap.Error(err)).Error("mongo.SingleResult.Decode failed in document.Document.GetOne")
		return err
	}

	return nil
}

func (im *impl) UpdateOne(ctx context.Context, name Name, query bson.M, doc interface{}) error {
	updater := bson.M{
		"$set": doc,
	}
	res := im.mongo.Database(dbName).Collection(name.String()).FindOneAndUpdate(ctx, query, updater)
	err := res.Err()
	if err == mongo.ErrNoDocuments {
		return ErrNotFound
	} else if err != nil {
		ctx.With(zap.Error(err)).Error("mongo.Collection.FindOneAndUpdate failed in document.Document.UpdateOne")
		return err
	}

	return nil
}

func (im *impl) DeleteOne(ctx context.Context, name Name, query bson.M) error {
	res := im.mongo.Database(dbName).Collection(name.String()).FindOneAndDelete(ctx, query)
	err := res.Err()
	if err == mongo.ErrNoDocuments {
		return ErrNotFound
	} else if err != nil {
		ctx.With(zap.Error(err)).Error("mongo.Collection.FindOneAndDelete failed in document.Document.DeleteOne")
		return err
	}

	return nil
}
