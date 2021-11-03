package repositories

import (
	"testing"

	"github.com/dominikus1993/dev-news-bot/src/core/model"
	"github.com/stretchr/testify/assert"
)

func TestFromArticles(t *testing.T) {
	articles := []model.Article{model.NewArticle("sdd", "dasas"), model.NewArticle("sdd", "dasas")}
	subject := fromArticles(articles)
	assert.NotNil(t, subject)
	assert.Len(t, subject, 2)
	mongoArt := subject[0].(mongoArticle)
	assert.Equal(t, "sdd", mongoArt.Title)
	assert.Equal(t, "dasas", mongoArt.Link)
}
