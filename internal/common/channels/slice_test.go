package channels

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func rangeInt(from, to int) []int {
	res := make([]int, 0, to-from)
	for i := from; i < to; i++ {
		res = append(res, i)
	}
	return res
}

// BenchmarkToSlice-4       1302279               897.8 ns/op           688 B/op         12 allocs/op
func TestToSlice(t *testing.T) {
	numbers := make(chan int, 10)
	go func() {
		for _, a := range rangeInt(0, 10) {
			numbers <- a
		}
		close(numbers)
	}()

	result := ToSlice(numbers)
	assert.Len(t, result, 10)
	assert.ElementsMatch(t, rangeInt(0, 10), result)
}

func BenchmarkToSlice(b *testing.B) {
	for n := 0; n < b.N; n++ {
		numbers := make(chan int, 10)
		go func() {
			for _, a := range rangeInt(0, 10) {
				numbers <- a
			}
			close(numbers)
		}()

		ToSlice(numbers)
	}
}
