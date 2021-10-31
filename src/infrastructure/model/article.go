package model

import (
	"time"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoArticle struct {
	Title     string             `bson:"_id"`
	Link      string             `bson:"Link"`
	Content   string             `bson:"Content"`
	CrawledAt primitive.DateTime `bson:"CrawledAt"`
}

func FromArticle(article *model.Article) *MongoArticle {
	return &MongoArticle{
		Title:     article.Title,
		Link:      article.Link,
		Content:   article.Content,
		CrawledAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
	}
}

func FromArticles(articles []model.Article) []interface{} {
	mongoArticles := make([]interface{}, len(articles))
	for i, article := range articles {
		mongoArticles[i] = *FromArticle(&article)
	}
	return mongoArticles
}
