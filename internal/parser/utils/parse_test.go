package utils

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	ctx := context.TODO()
	articles := []model.Article{
		model.NewArticle("title", "http://www.google.pl", "source"),
		model.NewArticle("title2", "http://www.google.pl", "source"),
	}
	result := Parse(ctx, func(ctx context.Context, stream chan<- model.Article) {
		for _, v := range articles {
			stream <- v
		}
	})
	subject := channels.ToSlice(result)
	assert.NotNil(t, subject)
	assert.NotEmpty(t, subject)
	assert.Len(t, subject, len(articles))
	for _, article := range subject {
		assert.True(t, article.IsValid())
	}
}
