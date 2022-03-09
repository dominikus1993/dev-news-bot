package mongo

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
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

func (r *mongoArticlesRepository) IsNew(ctx context.Context, article *model.Article) (bool, error) {
	opt := options.Count().SetLimit(1)
	res, err := r.countArticles(ctx, opt, bson.D{{"_id", article.Title}})
	if err != nil {
		return false, err
	}
	return res == 0, nil
}

func (r *mongoArticlesRepository) Save(ctx context.Context, articles []model.Article) error {
	col := r.getCollection()
	art := fromArticles(articles)
	_, err := col.InsertMany(ctx, art)
	if err != nil {
		return err
	}
	return nil
}

func (r *mongoArticlesRepository) countArticles(ctx context.Context, opts *options.CountOptions, filter interface{}) (int64, error) {
	col := r.getCollection()
	res, err := col.CountDocuments(ctx, filter, opts)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *mongoArticlesRepository) Read(ctx context.Context, params repositories.GetArticlesParams) (*repositories.Articles, error) {
	col := r.getCollection()
	opts := options.Find()
	if params.PageSize > 0 {
		opts.SetLimit(int64(params.PageSize))
	}
	if params.Page > 0 {
		opts.SetSkip(int64((params.Page - 1) * params.PageSize))
	}
	opts.SetSort(bson.D{{"CrawledAt", -1}})
	cur, err := col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var articles []model.Article
	for cur.Next(ctx) {
		var art mongoArticle
		if err := cur.Decode(&art); err != nil {
			return nil, err
		}
		articles = append(articles, toArticle(&art))
	}
	opt := options.Count()
	total, err := r.countArticles(ctx, opt, bson.D{})
	if err != nil {
		return nil, err
	}
	return repositories.NewArticles(articles, int(total)), nil
}
