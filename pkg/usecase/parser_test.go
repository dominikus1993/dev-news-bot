package usecase

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/notifications"
	"github.com/dominikus1993/dev-news-bot/pkg/providers"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
	"github.com/stretchr/testify/assert"
)

type fakeParser struct {
}

func (f *fakeParser) Parse(ctx context.Context) model.ArticlesStream {
	stream := make(chan model.Article)
	go func() {
		stream <- model.NewArticle("test", "http://dad", "reddit")
		close(stream)
	}()
	return stream
}

type fakeParser2 struct {
}

func (f *fakeParser2) Parse(ctx context.Context) model.ArticlesStream {
	stream := make(chan model.Article)
	go func() {
		defer close(stream)
		stream <- model.NewArticle("test", "http://dadsadad", "reddit")
	}()
	return stream
}

type fakeRepo struct {
	articles map[model.ArticleId]model.Article
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

func newFakeRepo() *fakeRepo {
	return &fakeRepo{articles: make(map[string]model.Article)}
}

type fakeNotifier struct {
	articles []model.Article
}

func (n *fakeNotifier) Notify(ctx context.Context, articles []model.Article) error {
	n.articles = append(n.articles, articles...)
	return nil
}

type fakeFilter struct {
}

func (f fakeFilter) Where(ctx context.Context, articles model.ArticlesStream) model.ArticlesStream {
	return articles
}

func TestParseArticlesAndSendItUseCase(t *testing.T) {
	provider := providers.NewArticlesProvider(newFakeRepo(), &fakeFilter{}, &fakeParser{}, &fakeParser2{})
	repo := newFakeRepo()
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster(notifier))

	ctx := context.Background()
	err := ucase.Execute(ctx, 10)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 2)
	assert.Len(t, notifier.articles, 2)
}

func TestParseArticlesAndSendItUseCaseWhenArticlesQuantityIsOne(t *testing.T) {
	provider := providers.NewArticlesProvider(newFakeRepo(), &fakeFilter{}, &fakeParser{}, &fakeParser2{})
	repo := newFakeRepo()
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster(notifier))

	ctx := context.Background()
	err := ucase.Execute(ctx, 1)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 1)
	assert.Len(t, notifier.articles, 1)
}

func TestParseArticlesWhenArticlesAleradyExistsInDb(t *testing.T) {
	provider := providers.NewArticlesProvider(newFakeRepo(), &fakeFilter{}, &fakeParser{}, &fakeParser2{})
	repo := newFakeRepo()
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster(notifier))

	ctx := context.Background()
	err := ucase.Execute(ctx, 1)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 1)
	assert.Len(t, notifier.articles, 1)
}

func TestParseArticlesAndSendItUseCaseWhenArticlesParserHasError(t *testing.T) {
	provider := providers.NewArticlesProvider(newFakeRepo(), &fakeFilter{}, &fakeParser{})
	repo := newFakeRepo()
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster(notifier))

	ctx := context.Background()
	err := ucase.Execute(ctx, 10)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 1)
	assert.Len(t, notifier.articles, 1)
}
