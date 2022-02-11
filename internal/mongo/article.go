package repositories

import (
	"time"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoArticle struct {
	Title     string             `bson:"_id"`
	Link      string             `bson:"Link"`
	Content   string             `bson:"Content"`
	CrawledAt primitive.DateTime `bson:"CrawledAt"`
}

func fromArticle(article *model.Article) *mongoArticle {
	return &mongoArticle{
		Title:     article.Title,
		Link:      article.Link,
		Content:   article.Content,
		CrawledAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
	}
}

func fromArticles(articles []model.Article) []interface{} {
	mongoArticles := make([]interface{}, len(articles))
	for i, article := range articles {
		mongoArticles[i] = *fromArticle(&article)
	}
	return mongoArticles
}
