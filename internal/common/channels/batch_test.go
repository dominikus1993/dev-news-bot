package channels

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	numbers := make(chan int, 10)
	go func() {
		for _, a := range rangeInt(1, 10) {
			numbers <- a
		}
		close(numbers)
	}()

	result := Batch(context.TODO(), numbers, 2)
	subject := ToSlice(result)
	assert.Len(t, subject, 5)
	assert.Equal(t, subject, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9}})
}

func BenchmarkBatcg(b *testing.B) {
	ctx := context.TODO()
	for n := 0; n < b.N; n++ {
		numbers := make(chan int, 10)
		go func() {
			for _, a := range rangeInt(1, 10) {
				numbers <- a
			}
			close(numbers)
		}()

		result := Batch(ctx, numbers, 2)
		ToSlice(result)
	}
}
