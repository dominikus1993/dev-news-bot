package mongo

import (
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoArticle struct {
	ID        string             `bson:"_id,omitempty"`
	Title     string             `bson:"Title"`
	Link      string             `bson:"Link"`
	Content   string             `bson:"Content"`
	CrawledAt primitive.DateTime `bson:"CrawledAt"`
}

func fromArticle(article *model.Article) *mongoArticle {
	return &mongoArticle{
		ID:        article.GetID(),
		Title:     article.GetTitle(),
		Link:      article.GetLink(),
		Content:   article.GetContent(),
		CrawledAt: primitive.NewDateTimeFromTime(article.GetCrawledAt()),
	}
}

func fromArticles(articles []model.Article) []interface{} {
	articles = model.UniqueArticlesArray(articles)
	mongoArticles := make([]interface{}, len(articles))

	for i, article := range articles {
		mongoArticles[i] = *fromArticle(&article)
	}
	return mongoArticles
}

func toArticle(article *mongoArticle) model.Article {
	return model.NewArticleWithContent(article.Title, article.Link, article.Content)
}
