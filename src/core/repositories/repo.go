package repositories

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
)

type IArticleRepository interface {
	Exists(ctx context.Context, article model.Article) (bool, error)
	Save(ctx context.Context, articles []model.Article) error
}
