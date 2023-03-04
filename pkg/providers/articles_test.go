package providers

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/stretchr/testify/assert"
)

type fakeParser struct {
}

func (f *fakeParser) Parse(ctx context.Context) model.ArticlesStream {
	stream := make(chan model.Article)
	go func() {
		stream <- model.NewArticle("test", "http://dad")
		close(stream)
	}()
	return stream
}

type fakeParser2 struct {
}

func (f *fakeParser2) Parse(ctx context.Context) model.ArticlesStream {
	stream := make(chan model.Article)
	go func() {
		stream <- model.NewArticle("test", "http://dadsadad")
		close(stream)
	}()
	return stream
}

type fakeRepo struct {
	articles map[model.ArticleId]model.Article
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{articles: make(map[string]model.Article)}
}

func (r *fakeRepo) FilterNew(ctx context.Context, stream model.ArticlesStream) model.ArticlesStream {
	result := make(chan model.Article)
	go func() {
		for article := range stream {
			if _, ok := r.articles[article.GetID()]; !ok {
				result <- article
			}
		}
		close(result)
	}()
	return result
}

func (r *fakeRepo) Read(ctx context.Context, params repositories.GetArticlesParams) (*repositories.Articles, error) {
	return nil, nil
}

func (r *fakeRepo) Save(ctx context.Context, articles []model.Article) error {
	for _, article := range articles {
		if _, ok := r.articles[article.GetID()]; !ok {
			r.articles[article.GetID()] = article
		}
	}
	return nil
}

type fakeFilter struct {
}

func (f fakeFilter) Where(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	return articles
}

func TestArticlesProvider(t *testing.T) {
	articlesProvider := NewArticlesProvider(newFakeRepo(), &fakeFilter{}, &fakeParser{}, &fakeParser2{})
	subject := channels.ToSlice(articlesProvider.Provide(context.Background()))
	assert.Len(t, subject, 2)
	assert.Equal(t, "test", subject[0].GetTitle())
	assert.Equal(t, "test", subject[1].GetTitle())
}

func TestArticlesProviderWhenArticlesAlreadyExistsInDb(t *testing.T) {
	repo := newFakeRepo()
	repo.Save(context.TODO(), []model.Article{model.NewArticle("test", "http://dad"), model.NewArticle("test", "http://dadsadad")})
	articlesProvider := NewArticlesProvider(repo, fakeFilter{}, &fakeParser{}, &fakeParser2{})
	subject := channels.ToSlice(articlesProvider.Provide(context.Background()))
	assert.Len(t, subject, 0)
}
