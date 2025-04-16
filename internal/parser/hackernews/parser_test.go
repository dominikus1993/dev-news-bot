package hackernews

import (
	"context"
	"testing"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/stretchr/testify/assert"
)

func TestHackerNews(t *testing.T) {
	const q = 10
	parser := NewHackerNewsArticleParser(q)
	result := parser.Parse(context.TODO())
	subject := channels.ToSlice(result)
	assert.NotNil(t, subject)
	assert.NotEmpty(t, subject)
	for _, article := range subject {
		assert.True(t, article.IsValid())
	}
}

func TestGetTopArticlesIds(t *testing.T) {
	ids, err := getTopArticlesIds()
	assert.Nil(t, err)
	assert.NotEmpty(t, ids)
}
