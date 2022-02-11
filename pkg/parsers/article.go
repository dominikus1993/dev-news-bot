package parsers

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
)

type ArticlesParser interface {
	Parse(ctx context.Context) ([]model.Article, error)
}
