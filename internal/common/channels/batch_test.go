package channels

import (
	"context"
	"testing"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/stretchr/testify/assert"
)

func rangeInt(from, to int) []int {
	res := make([]int, 0, to-from)
	for i := from; i < to; i++ {
		res = append(res, i)
	}
	return res
}

func TestBatch(t *testing.T) {
	numbers := make(chan int, 10)
	go func() {
		for _, a := range rangeInt(1, 10) {
			numbers <- a
		}
		close(numbers)
	}()

	result := Batch(context.TODO(), numbers, 2)
	subject := channels.ToSlice(result)
	assert.Len(t, subject, 5)
	assert.Equal(t, subject, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9}})
}

func TestRunInBatch(t *testing.T) {
	arr := rangeInt(1, 50)
	numbers := make(chan int, 10)
	go func() {
		for _, a := range arr {
			numbers <- a
		}
		close(numbers)
	}()

	result := RunBatchInPararell(context.TODO(), numbers, 10, 10, func(ctx context.Context, t []int) ([]int, error) {
		return t, nil
	})
	stream, _ := result.Read()
	subject := channels.ToSlice(stream)
	assert.Len(t, subject, len(arr))
	assert.ElementsMatch(t, arr, subject)
}

func BenchmarkBatcg(b *testing.B) {
	ctx := context.TODO()
	for b.Loop() {
		numbers := make(chan int, 10)
		go func() {
			for _, a := range rangeInt(1, 10) {
				numbers <- a
			}
			close(numbers)
		}()

		result := Batch(ctx, numbers, 2)
		channels.ToSlice(result)
	}
}
