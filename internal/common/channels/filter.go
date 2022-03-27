package channels

import "context"

func Filter[T any](ctx context.Context, source <-chan T, predicate func(ctx context.Context, article *T) bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case article, ok := <-source:
				if !ok {
					return
				}
				if predicate(ctx, &article) {
					out <- article
				}
			}
		}
	}()
	return out
}
