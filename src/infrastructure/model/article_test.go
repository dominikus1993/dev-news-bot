package model

import (
	"testing"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	"github.com/stretchr/testify/assert"
)

func TestFromArticles(t *testing.T) {
	articles := []model.Article{model.NewArticle("sdd", "dasas"), model.NewArticle("sdd", "dasas")}
	subject := FromArticles(articles)
	assert.NotNil(t, subject)
	assert.Len(t, subject, 2)
}
