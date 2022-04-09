package providers

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/internal/common/channels"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
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
		if a.GetTitle() == article.GetTitle() {
			return false, nil
		}
	}
	return true, nil
}

func (r *fakeRepo) Read(ctx context.Context, params repositories.GetArticlesParams) (*repositories.Articles, error) {
	return nil, nil
}

func TestArticlesProvider(t *testing.T) {
	articlesProvider := NewArticlesProvider(&fakeRepo{}, &fakeParser{}, &fakeParser2{})
	subject := channels.ToSlice(articlesProvider.Provide(context.Background()))
	assert.Len(t, subject, 2)
	assert.Equal(t, "test", subject[0].GetTitle())
	assert.Equal(t, "test", subject[1].GetTitle())
}

func TestArticlesProviderWhenArticlesAlreadyExistsInDb(t *testing.T) {
	articlesProvider := NewArticlesProvider(&fakeRepo{articles: []model.Article{model.NewArticle("test", "http://dad"), model.NewArticle("test", "http://dadsadad")}}, &fakeParser{}, &fakeParser2{})
	subject := channels.ToSlice(articlesProvider.Provide(context.Background()))
	assert.Len(t, subject, 0)
}
