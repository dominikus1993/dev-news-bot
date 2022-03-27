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
	provider := providers.NewArticlesProvider(&fakeRepo{}, &fakeParser{}, &fakeParser2{})
	repo := &fakeRepo{}
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster(notifier))

	ctx := context.Background()
	err := ucase.Execute(ctx, 10)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 2)
	assert.Len(t, notifier.articles, 2)
}

func TestParseArticlesAndSendItUseCaseWhenArticlesQuantityIsOne(t *testing.T) {
	provider := providers.NewArticlesProvider(&fakeRepo{}, &fakeParser{}, &fakeParser2{})
	repo := &fakeRepo{}
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster(notifier))

	ctx := context.Background()
	err := ucase.Execute(ctx, 1)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 1)
	assert.Len(t, notifier.articles, 1)
}

func TestParseArticlesWhenArticlesAleradyExistsInDb(t *testing.T) {
	provider := providers.NewArticlesProvider(&fakeRepo{}, &fakeParser{}, &fakeParser2{})
	repo := &fakeRepo{}
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster(notifier))

	ctx := context.Background()
	err := ucase.Execute(ctx, 1)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 1)
	assert.Len(t, notifier.articles, 1)
}

func TestParseArticlesAndSendItUseCaseWhenArticlesParserHasError(t *testing.T) {
	provider := providers.NewArticlesProvider(&fakeRepo{}, &fakeParser{})
	repo := &fakeRepo{}
	notifier := &fakeNotifier{}
	ucase := NewParseArticlesAndSendItUseCase(provider, repo, notifications.NewBroadcaster(notifier))

	ctx := context.Background()
	err := ucase.Execute(ctx, 10)
	assert.Nil(t, err)
	assert.Len(t, repo.articles, 1)
	assert.Len(t, notifier.articles, 1)
}
