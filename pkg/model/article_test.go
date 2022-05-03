package model

import (
	"testing"

	"github.com/dominikus1993/dev-news-bot/internal/common/channels"
	"github.com/stretchr/testify/assert"
)

func TestArticleValidationWitContent(t *testing.T) {
	article := NewArticleWithContent("test", "https://xd.pl", "content")
	subject := article.IsValid()
	assert.True(t, subject)
}

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
	articles := make(chan Article, 10)
	for _, a := range []Article{NewArticle("x", "2"), NewArticle("d", "1"), NewArticle("xd", "37")} {
		articles <- a
	}
	close(articles)
	randomArticles := TakeRandomArticles(articles, 0)
	assert.Len(t, randomArticles, 0)
}

func TestGetRandomArticlesWhenTakeIsGreaterThanLenOfArticlesArray(t *testing.T) {
	articles := make(chan Article, 10)
	for _, a := range []Article{NewArticle("x", "2"), NewArticle("d", "1"), NewArticle("xd", "37")} {
		articles <- a
	}
	close(articles)
	randomArticles := TakeRandomArticles(articles, 5)
	assert.Len(t, randomArticles, 3)
}

func TestGetRandomArticlesWhenTakeIsSmallerThanLenOfArticlesArray(t *testing.T) {
	articles := make(chan Article, 10)
	for _, a := range []Article{NewArticle("x", "2"), NewArticle("d", "1"), NewArticle("xd", "37")} {
		articles <- a
	}
	close(articles)
	randomArticles := TakeRandomArticles(articles, 2)
	assert.Len(t, randomArticles, 2)
}

//96d86055-c7c2-570e-91a7-06567337af70

func TestGetUniqueArticlesFromStream(t *testing.T) {
	articles := make(chan Article, 10)
	for _, a := range []Article{NewArticle("x", "2"), NewArticle("d", "1"), NewArticle("xd", "37"), NewArticle("x", "2"), NewArticle("x", "3"), NewArticle("xd", "37")} {
		articles <- a
	}
	close(articles)
	uniqueArticles := UniqueArticles(articles)
	subject := channels.ToSlice(uniqueArticles)
	assert.Len(t, subject, 4)
}
