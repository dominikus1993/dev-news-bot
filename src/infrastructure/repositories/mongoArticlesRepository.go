package repositories

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoArticlesRepository struct {
	client *MongoClient
	db     *mongo.Database
}

func NewMongoArticlesRepository(client *MongoClient, database string) *mongoArticlesRepository {
	return &mongoArticlesRepository{client: client, db: client.mongo.Database(database)}
}

func (c *mongoArticlesRepository) getCollection() *mongo.Collection {
	return c.db.Collection("articles")
}

func (r *mongoArticlesRepository) Exists(ctx context.Context, article model.Article) (bool, error) {
	col := r.getCollection()
	opts := options.Count().SetLimit(1)
	res, err := col.CountDocuments(ctx, bson.D{{"_id", article.Title}}, opts)
	if err != nil {
		return false, err
	}
	return res > 0, nil
}
