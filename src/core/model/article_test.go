package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleValidationWhenUrlIsInCorrect(t *testing.T) {
	article := NewArticle("test", "notlink")
	subject := article.IsValid()
	assert.False(t, subject)
}
