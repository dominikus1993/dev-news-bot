package filters

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
)

type ArticlesFilter interface {
	Where(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream
}
