package mongo

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	mongo      *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewClient(ctx context.Context, connectionString, database string) (*MongoClient, error) {

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}
	db := client.Database(database)
	collection := db.Collection("articles")

	return &MongoClient{mongo: client, db: db, collection: collection}, nil
}

func (c *MongoClient) Close(ctx context.Context) {
	if err := c.mongo.Disconnect(ctx); err != nil {
		log.WithError(err).Error("Error when trying disconnect from mongo")
	}
}
