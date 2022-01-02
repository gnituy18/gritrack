package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"gritter/pkg/log"
)

var (
	mongoClient *mongo.Client
)

func GetMongo() (*mongo.Client, error) {
	if mongoClient != nil {
		return mongoClient, nil
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Global().With(zap.Error(err)).Error("mongo.NewClient failed in db.GetMongo")
		return nil, err
	}

	if err := client.Connect(context.Background()); err != nil {
		panic(err)
	}

	mongoClient = client

	return client, nil
}

func MustGetMongo() *mongo.Client {
	client, err := GetMongo()
	if err != nil {
		log.Global().With(zap.Error(err)).Error("db.GetMongo failed in db.MustGetMongo")
		panic(err)
	}

	return client
}
