package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIntWhenEmpty(t *testing.T) {
	str := ""
	subject := ParseInt(str, 1)
	assert.Equal(t, 1, subject)
}

func TestParseIntWhenIncorrectString(t *testing.T) {
	str := "ddsdds"
	subject := ParseInt(str, 1)
	assert.Equal(t, 1, subject)
}

func TestParseInt(t *testing.T) {
	str := "3"
	subject := ParseInt(str, 1)
	assert.Equal(t, 3, subject)
}

func TestParseIntWhenNegative(t *testing.T) {
	str := "-3"
	subject := ParseInt(str, 1)
	assert.Equal(t, -3, subject)
}
