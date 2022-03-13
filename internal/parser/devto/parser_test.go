package devto

import (
	"context"
	"testing"

	"github.com/dominikus1993/dev-news-bot/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestDevToParser(t *testing.T) {
	parser := NewDevToParser([]string{"golang", "dotnet"})
	stream := parser.Parse(context.TODO())
	subject := model.ToArticlesArray(stream)
	assert.NotEmpty(t, subject)
}
