package channels

import (
	"context"
	"sync"
)

func FanIn[T any](ctx context.Context, stream ...<-chan T) chan T {
	var wg sync.WaitGroup
	out := make(chan T)
	output := func(c <-chan T) {
		defer wg.Done()
		for v := range c {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}
	wg.Add(len(stream))
	for _, c := range stream {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
