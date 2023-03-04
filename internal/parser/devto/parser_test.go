package devto

import (
	"context"
	"testing"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/stretchr/testify/assert"
)

func TestDevToParser(t *testing.T) {
	parser := NewDevToParser([]string{"golang", "dotnet"})
	stream := parser.Parse(context.TODO())
	subject := channels.ToSlice(stream)
	assert.NotEmpty(t, subject)
}
