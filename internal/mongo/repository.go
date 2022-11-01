package mongo

import (
	"context"
	"errors"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoArticlesRepository struct {
	client *MongoClient
}

func NewMongoArticlesRepository(client *MongoClient) *mongoArticlesRepository {
	return &mongoArticlesRepository{client: client}
}

func (r *mongoArticlesRepository) IsNew(ctx context.Context, article model.Article) (bool, error) {
	col := r.client.collection
	opts := options.FindOne()
	res := col.FindOne(ctx, bson.D{{Key: "_id", Value: article.GetID()}}, opts)
	notexists := res.Err() == mongo.ErrNoDocuments
	return notexists, nil
}

func (r *mongoArticlesRepository) Save(ctx context.Context, articles []model.Article) error {
	if len(articles) == 0 {
		return errors.New("no articles to save")
	}
	session, err := r.client.mongo.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	err = session.StartTransaction()
	if err != nil {
		return err
	}
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		art := fromArticles(articles)

		_, err = r.client.collection.InsertMany(ctx, art)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		session.AbortTransaction(ctx)
		return err
	}
	return session.CommitTransaction(ctx)
}

func (r *mongoArticlesRepository) Read(ctx context.Context, params repositories.GetArticlesParams) (*repositories.Articles, error) {
	col := r.client.collection

	opts := options.Find()
	if params.PageSize > 0 {
		opts.SetLimit(int64(params.PageSize))
	}
	if params.Page > 0 {
		opts.SetSkip(int64((params.Page - 1) * params.PageSize))
	}
	opts.SetSort(bson.D{{Key: "CrawledAt", Value: -1}})
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
	total, err := col.CountDocuments(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	}
	return repositories.NewArticles(articles, int(total)), nil
}
