package devto

import (
	"testing"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/stretchr/testify/assert"
)

func TestDevToParser(t *testing.T) {
	parser := NewDevToParser([]string{"golang", "dotnet"})
	stream := parser.Parse(t.Context())
	subject := channels.ToSlice(stream)
	assert.NotEmpty(t, subject)
}
