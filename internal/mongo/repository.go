package mongo

import (
	"context"
	"errors"

	"github.com/dominikus1993/dev-news-bot/internal/common/channels"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var projectionStage = bson.D{{"$project", bson.D{{"_id", "$_id"}}}}

type mongoArticlesRepository struct {
	client *MongoClient
}

func NewMongoArticlesRepository(client *MongoClient) *mongoArticlesRepository {
	return &mongoArticlesRepository{client: client}
}

func (r *mongoArticlesRepository) getArticlesCollection() *mongo.Collection {
	return r.client.collection
}

func getArticlesIds(articles []model.Article) []model.ArticleId {
	result := make([]model.ArticleId, len(articles))

	for i, article := range articles {
		result[i] = article.GetID()
	}

	return result
}

type mongoArticleId struct {
	ID string `bson:"_id,omitempty"`
}

func contains(ids []mongoArticleId, el model.ArticleId) bool {
	length := len(ids)
	if length == 0 {
		return false
	}
	for i := 0; i < length; i++ {
		if ids[i].ID == el {
			return true
		}
	}
	return false
}

type articleExistence struct {
	article    model.Article
	existsInDb bool
}

func toMap(ids []mongoArticleId, articles []model.Article) map[model.ArticleId]articleExistence {
	result := make(map[model.ArticleId]articleExistence)
	for _, v := range articles {
		var id model.ArticleId = v.GetID()
		exists := contains(ids, id)
		result[id] = articleExistence{article: v, existsInDb: exists}
	}
	return result
}

// [{ "$match" : { "_id" : { "$in" : ["xDDD", "534ee62d-bd5d-5fa8-b384-73a70d8503b6"] } } }, { "$project" : { "_id" : "$_id" } }]
func (r *mongoArticlesRepository) checkArticleExistence(ctx context.Context, articles ...model.Article) (map[model.ArticleId]articleExistence, error) {
	result := make(map[model.ArticleId]articleExistence)
	if len(articles) == 0 {
		return result, nil
	}
	matchStage := bson.D{{"$match", bson.D{{"_id", bson.D{{"$in", getArticlesIds(articles)}}}}}}

	cursor, err := r.client.collection.Aggregate(ctx, mongo.Pipeline{matchStage, projectionStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []mongoArticleId
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return toMap(results, articles), nil
}

func (r *mongoArticlesRepository) filterOldArticles(ctx context.Context, articles []model.Article) ([]model.Article, error) {
	result := make([]model.Article, 0)
	existenceMap, err := r.checkArticleExistence(ctx, articles...)
	if err != nil {
		return nil, err
	}
	for _, article := range articles {
		if x, found := existenceMap[article.GetID()]; found && !x.existsInDb {
			result = append(result, x.article)
		}
	}
	return result, nil
}

func (r *mongoArticlesRepository) FilterNew(ctx context.Context, stream model.ArticlesStream) model.ArticlesStream {
	result := make(chan model.Article)
	go func() {
		batch := channels.Batch(ctx, stream, 10)
		for articles := range batch {
			newArticles, err := r.filterOldArticles(ctx, articles)
			if err != nil {
				log.WithError(err).Errorln("can't check mongo articles existence")
				break
			}
			for i := 0; i < len(newArticles); i++ {
				result <- newArticles[i]
			}
		}
		close(result)
	}()
	return result
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

		_, err = r.client.collection.BulkWrite(ctx, art)
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
		opts = opts.SetLimit(int64(params.PageSize))
	}
	if params.Page > 0 {
		opts = opts.SetSkip(int64((params.Page - 1) * params.PageSize))
	}
	opts = opts.SetSort(bson.D{{Key: "CrawledAt", Value: -1}})
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
