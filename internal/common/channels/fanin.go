package channels

import (
	"context"
	"sync"
)

func FanIn[T any](ctx context.Context, stream ...<-chan T) chan T {
	var wg sync.WaitGroup
	out := make(chan T)
	output := func(c <-chan T) {
		for v := range c {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}
	for _, c := range stream {
		wg.Go(func() {
			output(c)
		})
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()
	return out
}
