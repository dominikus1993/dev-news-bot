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

func TestArticleValidationWhenTileIsEmpty(t *testing.T) {
	article := NewArticle("", "notlink")
	subject := article.IsValid()
	assert.False(t, subject)
}

func TestArticleValidationWhenLinkIsEmpty(t *testing.T) {
	article := NewArticle("asdddddddd", "")
	subject := article.IsValid()
	assert.False(t, subject)
}

func TestArticleValidationWhenUrlIsCorrect(t *testing.T) {
	article := NewArticle("asdddddddd", "https://scienceintegritydigest.com/2020/12/20/paper-about-herbalife-related-patient-death-removed-after-company-threatens-to-sue-the-journal/")
	subject := article.IsValid()
	assert.True(t, subject)
}
