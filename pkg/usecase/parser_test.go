package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/dev-news-bot/pkg/notifications"
	"github.com/dominikus1993/dev-news-bot/pkg/parsers"
	"github.com/dominikus1993/dev-news-bot/pkg/providers"
	"github.com/dominikus1993/dev-news-bot/pkg/repositories"
	"github.com/stretchr/testify/assert"
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

type fakeErrorParser struct {
}

func (f *fakeErrorParser) Parse(ctx context.Context) ([]model.Article, error) {
	return nil, errors.New("xDDD")
}

type fakeRepo struct {
	articles []model.Article
}

func (r *fakeRepo) IsNew(ctx context.Context, article *model.Article) (bool, error) {
	for _, a := range r.articles {
		if a.Title == article.Title {
			return false, nil
		}
	}
	return true, nil
}

func (r *fakeRepo) Read(ctx context.Context, params repositories.GetArticlesParams) (*repositories.Articles, error) {
	return nil, nil
}

func (r *fakeRepo) Save(ctx context.Context, articles []model.Article) error {
	r.articles = append(r.articles, articles...)
	return nil
}

type fakeNotifier struct {
	articles []model.Article
}

func (n *fakeNotifier) Notify(ctx context.Context, articles []model.Article) error {
	n.articles = append(n.articles, articles...)
	return nil
}

func TestParseArticlesAndSendItUseCase(t *testing.T) {
	provider := providers.NewArticlesProvider([]parsers.ArticlesParser{
		&fakeParser{},
		&fakeParser2{},
	})
	repo := &fakeRepo{}
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster([]notifications.Notifier{
		notifier,
	}))

	ctx := context.Background()
	err := ucase.Execute(ctx, 10)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 2)
	assert.Len(t, notifier.articles, 2)
}

func TestParseArticlesAndSendItUseCaseWhenArticlesQuantityIsOne(t *testing.T) {
	provider := providers.NewArticlesProvider([]parsers.ArticlesParser{
		&fakeParser{},
		&fakeParser2{},
	})
	repo := &fakeRepo{}
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster([]notifications.Notifier{
		notifier,
	}))

	ctx := context.Background()
	err := ucase.Execute(ctx, 1)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 1)
	assert.Len(t, notifier.articles, 1)
}

func TestParseArticlesAndSendItUseCaseWhenArticlesParserHasError(t *testing.T) {
	provider := providers.NewArticlesProvider([]parsers.ArticlesParser{
		&fakeParser{},
		&fakeErrorParser{},
	})
	repo := &fakeRepo{}
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster([]notifications.Notifier{
		notifier,
	}))

	ctx := context.Background()
	err := ucase.Execute(ctx, 10)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 1)
	assert.Len(t, notifier.articles, 1)
}
