package channels

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFanIn(t *testing.T) {
	numbers := make(chan int, 10)
	numbers2 := make(chan int, 10)
	go func() {
		for _, a := range rangeInt(1, 10) {
			numbers <- a
		}
		close(numbers)
	}()

	go func() {
		for _, a := range rangeInt(10, 20) {
			numbers2 <- a
		}
		close(numbers2)
	}()

	result := FanIn(context.TODO(), numbers, numbers2)
	subject := ToSlice(result)
	assert.Len(t, subject, 19)
	assert.ElementsMatch(t, rangeInt(1, 20), subject)
}

func BenchmarkFanIn(b *testing.B) {
	ctx := context.TODO()
	for n := 0; n < b.N; n++ {
		numbers := make(chan int, 10)
		numbers2 := make(chan int, 10)
		go func() {
			for _, a := range rangeInt(1, 10) {
				numbers <- a
			}
			close(numbers)
		}()

		go func() {
			for _, a := range rangeInt(10, 20) {
				numbers2 <- a
			}
			close(numbers2)
		}()

		result := FanIn(ctx, numbers, numbers2)
		ToSlice(result)
	}
}
