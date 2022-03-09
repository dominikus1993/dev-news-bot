package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	pageSize int
	total    int
	expected int
}

func TestCountNumberOfPages(t *testing.T) {
	data := []testData{
		{pageSize: 10, total: 99, expected: 10},
		{pageSize: 10, total: 100, expected: 10},
		{pageSize: 10, total: 101, expected: 11},
		{pageSize: 1, total: 99, expected: 99},
		{pageSize: 20, total: 10, expected: 1},
	}

	for _, element := range data {
		pages := countNumberOfPages(element.total, element.pageSize)
		assert.Equal(t, element.expected, pages)
	}
}
