package channels

import "context"

func Filter[T any](ctx context.Context, source <-chan T, predicate func(ctx context.Context, article *T) bool) <-chan T {
	result := make(chan T, 100)
	go func() {
		for data := range source {
			if predicate(ctx, &data) {
				result <- data
			}
		}
		close(result)
	}()
	return result
}
