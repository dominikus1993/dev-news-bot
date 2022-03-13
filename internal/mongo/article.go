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
		ID:        article.ID,
		Title:     article.Title,
		Link:      article.Link,
		Content:   article.Content,
		CrawledAt: primitive.NewDateTimeFromTime(article.CrawledAt),
	}
}

func fromArticles(articles []model.Article) []interface{} {
	mongoArticles := make([]interface{}, len(articles))
	for i, article := range articles {
		mongoArticles[i] = *fromArticle(&article)
	}
	return mongoArticles
}

func toArticle(article *mongoArticle) model.Article {
	return model.Article{
		ID:        article.ID,
		Title:     article.Title,
		Link:      article.Link,
		Content:   article.Content,
		CrawledAt: article.CrawledAt.Time(),
	}
}
