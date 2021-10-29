package providers

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
)

type fakeParser struct {
}

func (f *fakeParser) Parse(ctx context.Context) ([]model.Article, error) {
	return []model.Article{model.NewArticle("test", "http://dad")}, nil
}

type fakeParser2 struct {
}

func (f *fakeParser2) Parse(ctx context.Context) ([]model.Article, error) {
	return []model.Article{model.NewArticle("test", "http://dadsadad")}, nil
}
