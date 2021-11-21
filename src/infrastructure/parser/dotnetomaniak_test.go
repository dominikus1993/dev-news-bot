package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDotnetomaniakParser(t *testing.T) {
	parser := NewDotnetoManiakParser()
	result, err := parser.Parse(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result)
	for _, article := range result {
		assert.True(t, article.IsValid())
	}
}
