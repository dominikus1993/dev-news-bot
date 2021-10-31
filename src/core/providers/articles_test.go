package providers

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	"github.com/dominikus1993/dev-news-bot/src/core/parsers"
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

func TestArticlesProvider(t *testing.T) {
	articlesProvider := NewArticlesProvider([]parsers.ArticlesParser{&fakeParser{}, &fakeParser2{}})
	subject, err := articlesProvider.Provide(context.Background())
	assert.Nil(t, err)
	assert.Len(t, subject, 2)
	assert.Equal(t, "test", subject[0].Title)
	assert.Equal(t, "test", subject[1].Title)
}