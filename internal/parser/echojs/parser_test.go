package echojs

import (
	"context"
	"testing"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/stretchr/testify/assert"
)

func TestEchoJs(t *testing.T) {
	parser := NewEechoJsParser()
	result := parser.Parse(context.TODO())
	subject := channels.ToSlice(result)
	assert.NotNil(t, subject)
	assert.NotEmpty(t, subject)
	for _, article := range subject {
		assert.True(t, article.IsValid())
	}
}
