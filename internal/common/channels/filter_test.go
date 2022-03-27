package channels

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	numbers := make(chan int, 10)
	go func() {
		for _, a := range rangeInt(1, 10) {
			numbers <- a
		}
		close(numbers)
	}()

	result := Filter(context.TODO(), numbers, func(ctx context.Context, element int) bool { return element%2 == 0 })
	subject := ToSlice(result)
	assert.Len(t, subject, 4)
	assert.ElementsMatch(t, []int{2, 4, 6, 8}, subject)
}

func BenchmarkFilter(b *testing.B) {
	ctx := context.TODO()
	for n := 0; n < b.N; n++ {
		numbers := make(chan int, 10)
		go func() {
			for _, a := range rangeInt(0, 10) {
				numbers <- a
			}
			close(numbers)
		}()

		ToSlice(Filter(ctx, numbers, func(ctx context.Context, element int) bool { return element%2 == 0 }))
	}
}
