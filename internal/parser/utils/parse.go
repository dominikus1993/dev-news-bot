package utils

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
)

func Parse(ctx context.Context, action func(ctx context.Context, stream chan<- model.Article)) model.ArticlesStream {
	result := make(chan model.Article, 20)
	go func() {
		defer close(result)
		action(ctx, result)
	}()
	return result
}
