package language

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/internal/common/channels"
	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/stretchr/testify/assert"
)

func fromSlice[T any](s []T) chan T {
	res := make(chan T)
	go func() {
		for _, v := range s {
			res <- v
		}
		close(res)
	}()
	return res
}

func TestGetArticleTitleAndContenWhenArticleHasNoContent(t *testing.T) {
	article := model.NewArticle("77% of all developers are involved in DevOps", "https://dev.to/slashdatahq/77-of-all-developers-are-involved-in-devops-hao")
	subject := getArticleTitleAndContent(article)
	assert.NotNil(t, subject)
	assert.NotEmpty(t, subject)
	assert.Equal(t, article.GetTitle(), subject)
}

func TestGetArticleTitleAndContenWhenArticleHasContent(t *testing.T) {
	article := model.NewArticleWithContent("77% of all developers are involved in DevOps", "https://dev.to/slashdatahq/77-of-all-developers-are-involved-in-devops-hao", "Some Content")
	subject := getArticleTitleAndContent(article)
	assert.NotNil(t, subject)
	assert.NotEmpty(t, subject)
	assert.Equal(t, "77% of all developers are involved in DevOps Some Content", subject)
}

func TestLanguageFilter(t *testing.T) {
	filter := NewLanguageFilter()
	articles := []model.Article{
		model.NewArticle("77% of all developers are involved in DevOps", "https://dev.to/slashdatahq/77-of-all-developers-are-involved-in-devops-hao"),
		model.NewArticle("Nauka w Pracy. Czy Programista Powinien Mieć Czas Na Naukę w Pracy? - Modest Programmer", "https://dotnetomaniak.pl/Nauka-w-Pracy-Czy-Programista-Powinien-Miec-Czas-Na-Nauke-w-Pracy-Modest-Programmer"),
		model.NewArticle("Executando Tasks em paralelo em aplicações .Net com o SemaphoreSlim", "https://dev.to/marcosbelorio/executando-tasks-em-paralelo-em-aplicacoes-net-com-o-semaphoreslim-4e9n")}
	result := filter.Where(context.TODO(), fromSlice(articles))
	subject := channels.ToSlice(result)
	assert.NotNil(t, subject)
	assert.NotEmpty(t, subject)
	assert.Len(t, subject, 2)
	for _, article := range subject {
		assert.True(t, article.IsValid())
	}
}
