package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUIDFromArticleTitle(t *testing.T) {
	article := NewArticleWithContent("test", "https://xd.pl", "content")
	subject, err := GenerateArticleId(article)
	assert.NoError(t, err)
	assert.Equal(t, "51614d0f-e50f-5728-931a-923553fca3d7", subject)
}

func TestUUIDWhenKeysAreEmpty(t *testing.T) {
	_, err := GenerateId()
	assert.Error(t, err)
}
