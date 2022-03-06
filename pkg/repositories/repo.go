package repositories

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
)

type GetArticlesParams struct {
	Page     int
	PageSize int
}

type ArticlesReader interface {
	IsNew(ctx context.Context, article *model.Article) (bool, error)
	Read(ctx context.Context, params GetArticlesParams) ([]model.Article, error)
}

type ArticlesWriter interface {
	Save(ctx context.Context, articles []model.Article) error
}

type ArticleRepository interface {
	ArticlesReader
	ArticlesWriter
}
