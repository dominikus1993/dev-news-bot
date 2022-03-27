package hackernews

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/internal/common/channels"
	"github.com/stretchr/testify/assert"
)

func TestHackerNews(t *testing.T) {
	parser := NewHackerNewsArticleParser()
	result := parser.Parse(context.TODO())
	subject := channels.ToSlice(result)
	assert.NotNil(t, subject)
	assert.NotEmpty(t, subject)
	for _, article := range subject {
		assert.True(t, article.IsValid())
	}
}
