package repositories

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	mongo *mongo.Client
}

func NewClient(connectionString string, ctx context.Context) *MongoClient {

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.WithContext(ctx).WithField("ConnectionString", connectionString).WithError(err).Fatal("Error when trying connect to mongo")
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.WithContext(ctx).WithField("ConnectionString", connectionString).WithError(err).Fatal("Error when trying ping mongo")
	}

	return &MongoClient{mongo: client}
}

func (c *MongoClient) Close(ctx context.Context) {
	if err := c.mongo.Disconnect(ctx); err != nil {
		log.WithError(err).Error("Error when trying disconnect from mongo")
	}
}
