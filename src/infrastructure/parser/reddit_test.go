package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedditParser(t *testing.T) {
	parser := NewRedditParser([]string{"golang", "dotnet"})
	subject, err := parser.Parse(context.TODO())
	assert.Nil(t, err)
	assert.NotEmpty(t, subject)
}
