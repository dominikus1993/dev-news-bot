package mongo

import (
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestFromArticles(t *testing.T) {
	articles := []model.Article{model.NewArticle("sdd", "dasas", "reddit"), model.NewArticle("sdd", "dasas", "reddit")}
	subject := fromArticles(articles)
	assert.NotNil(t, subject)
	assert.Len(t, subject, 1)
}
