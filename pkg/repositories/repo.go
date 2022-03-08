package repositories

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
)

type GetArticlesParams struct {
	Page     int
	PageSize int
}

type Articles struct {
	Articles []model.Article
	Total    int
}

func NewArticles(articles []model.Article, total int) *Articles {
	return &Articles{Articles: articles, Total: total}
}

type ArticlesReader interface {
	IsNew(ctx context.Context, article *model.Article) (bool, error)
	Read(ctx context.Context, params GetArticlesParams) (*Articles, error)
}

type ArticlesWriter interface {
	Save(ctx context.Context, articles []model.Article) error
}

type ArticleRepository interface {
	ArticlesReader
	ArticlesWriter
}
