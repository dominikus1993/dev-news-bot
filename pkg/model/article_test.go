package model

import (
	"testing"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/stretchr/testify/assert"
)

func TestArticleValidationWitContent(t *testing.T) {
	article := NewArticleWithContent("test", "https://xd.pl", "content", "reddit")
	subject := article.IsValid()
	assert.True(t, subject)
}

func TestArticleValidationWhenUrlIsEmptyt(t *testing.T) {
	article := NewArticle("test", "", "reddit")
	subject := article.IsValid()
	assert.False(t, subject)
}

func TestArticleValidationWhenUrlIsInCorrect(t *testing.T) {
	article := NewArticle("test", "notlink", "reddit")
	subject := article.IsValid()
	assert.False(t, subject)
}

func TestArticleValidationWhenTileIsEmpty(t *testing.T) {
	article := NewArticle("", "notlink", "reddit")
	subject := article.IsValid()
	assert.False(t, subject)
}

func TestArticleValidationWhenLinkIsEmpty(t *testing.T) {
	article := NewArticle("asdddddddd", "", "reddit")
	subject := article.IsValid()
	assert.False(t, subject)
}

func TestArticleValidationWhenUrlIsCorrect(t *testing.T) {
	article := NewArticle("asdddddddd", "https://scienceintegritydigest.com/2020/12/20/paper-about-herbalife-related-patient-death-removed-after-company-threatens-to-sue-the-journal/", "reddit")
	subject := article.IsValid()
	assert.True(t, subject)
}

func TestGetRandomArticlesWhenTakeIsZero(t *testing.T) {
	articles := make(chan Article, 10)
	for _, a := range []Article{NewArticle("x", "2", "reddit"), NewArticle("d", "1", "reddit"), NewArticle("xd", "37", "reddit")} {
		articles <- a
	}
	close(articles)
	randomArticles := TakeRandomArticles(articles, 0)
	assert.Len(t, randomArticles, 0)
}

func TestGetRandomArticlesWhenTakeIsGreaterThanLenOfArticlesArray(t *testing.T) {
	articles := make(chan Article, 10)
	for _, a := range []Article{NewArticle("x", "2", "reddit"), NewArticle("d", "1", "reddit"), NewArticle("xd", "37", "reddit")} {
		articles <- a
	}
	close(articles)
	randomArticles := TakeRandomArticles(articles, 5)
	assert.Len(t, randomArticles, 3)
}

func TestGetRandomArticlesWhenTakeIsSmallerThanLenOfArticlesArray(t *testing.T) {
	articles := make(chan Article, 10)
	for _, a := range []Article{NewArticle("x", "2", "reddit"), NewArticle("d", "1", "reddit"), NewArticle("xd", "37", "reddit")} {
		articles <- a
	}
	close(articles)
	randomArticles := TakeRandomArticles(articles, 2)
	assert.Len(t, randomArticles, 2)
}

//96d86055-c7c2-570e-91a7-06567337af70

func TestGetUniqueArticlesFromStream(t *testing.T) {
	articles := make(chan Article, 10)
	go func() {
		for _, a := range []Article{NewArticle("x", "2", "reddit"), NewArticle("d", "1", "reddit"), NewArticle("xd", "37", "reddit"), NewArticle("x", "2", "reddit"), NewArticle("x", "3", "reddit"), NewArticle("xd", "37", "reddit")} {
			articles <- a
		}
		close(articles)
	}()
	uniqueArticles := UniqueArticles(articles)
	subject := channels.ToSlice(uniqueArticles)
	assert.Len(t, subject, 4)
}

func TestTakeRandomToSlice(t *testing.T) {
	testArticles := []Article{NewArticle("x", "2", "reddit"), NewArticle("d", "1", "reddit"), NewArticle("xd", "37", "reddit"), NewArticle("x", "2", "reddit"), NewArticle("x", "3", "reddit"), NewArticle("xd", "37", "reddit")}

	subject := takeRandomToSlice(channels.FromSlice(testArticles), 3)
	assert.Len(t, subject, 3)
}

func TestTakeRandomToSliceWhenTakeIsGreaterThanLengthOfSourceSlice(t *testing.T) {
	testArticles := []Article{NewArticle("x", "2", "reddit")}

	subject := takeRandomToSlice(channels.FromSlice(testArticles), 3)
	assert.Len(t, subject, 1)
}
