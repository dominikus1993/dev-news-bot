package mongo

import (
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoArticle struct {
	ID        string             `bson:"_id,omitempty"`
	Title     string             `bson:"Title"`
	Link      string             `bson:"Link"`
	Content   string             `bson:"Content"`
	Source    string             `bson:"Source"`
	CrawledAt primitive.DateTime `bson:"CrawledAt"`
}

func fromArticle(article *model.Article) mongo.WriteModel {
	return mongo.NewInsertOneModel().SetDocument(mongoArticle{
		ID:        article.GetID(),
		Title:     article.GetTitle(),
		Link:      article.GetLink(),
		Content:   article.GetContent(),
		Source:    article.GetSource(),
		CrawledAt: primitive.NewDateTimeFromTime(article.GetCrawledAt()),
	})
}

func fromArticles(articles []model.Article) []mongo.WriteModel {
	articles = model.UniqueArticlesArray(articles)
	mongoArticles := make([]mongo.WriteModel, len(articles))

	for i, article := range articles {
		mongoArticles[i] = fromArticle(&article)
	}
	return mongoArticles
}

func toArticle(article *mongoArticle) model.Article {
	return model.NewArticleWithContent(article.Title, article.Link, article.Content, article.Source)
}
