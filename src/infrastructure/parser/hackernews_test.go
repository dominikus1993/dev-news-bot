package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHackerNews(t *testing.T) {
	parser := NewHackerNewsArticleParser()
	result, err := parser.Parse(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result)
}
