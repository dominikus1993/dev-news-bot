package providers

import (
	"context"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
)

type ArticlesProvider interface {
	Provide(ctx context.Context) ([]model.Article, error)
}

type fakeArticlesProvider struct {
}

func NewFakeArticlesProvider() *fakeArticlesProvider {
	return &fakeArticlesProvider{}
}

func (f *fakeArticlesProvider) Provide(ctx context.Context) ([]model.Article, error) {
	return []model.Article{
		{
			Title: "Fake article 1",
			Link:  "https://fake.com/article1",
		},
		{
			Title: "Fake article 2",
			Link:  "https://fake.com/article2",
		},
	}, nil
}
