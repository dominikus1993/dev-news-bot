package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleValidationWhenUrlIsEmptyt(t *testing.T) {
	article := NewArticle("test", "")
	subject := article.IsValid()
	assert.False(t, subject)
}

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

func TestGetRandomArticlesWhenTakeIsZero(t *testing.T) {
	articles := []Article{NewArticle("x", "2"), NewArticle("d", "1"), NewArticle("xd", "37")}
	randomArticles := GetRandomArticles(articles, 0)
	assert.Len(t, randomArticles, 0)
}

func TestGetRandomArticlesWhenTakeIsGreaterThanLenOfArticlesArray(t *testing.T) {
	articles := []Article{NewArticle("x", "2"), NewArticle("d", "1"), NewArticle("xd", "37")}
	randomArticles := GetRandomArticles(articles, 5)
	assert.Len(t, randomArticles, len(articles))
}

func TestGetRandomArticlesWhenTakeIsSmallerThanLenOfArticlesArray(t *testing.T) {
	articles := []Article{NewArticle("x", "2"), NewArticle("d", "1"), NewArticle("xd", "37")}
	randomArticles := GetRandomArticles(articles, 2)
	assert.Len(t, randomArticles, 2)
}