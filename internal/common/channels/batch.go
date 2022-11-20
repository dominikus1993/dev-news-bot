package channels

import (
	"context"
)

func sendToChannel[T any](batches chan []T, batch []T) {
	if len(batch) > 0 {
		batches <- batch
	}
}

func Batch[T any](ctx context.Context, values <-chan T, maxItems int) chan []T {
	batches := make(chan []T)
	go func() {
		defer close(batches)
		var batch []T
		for value := range values {
			batch = append(batch, value)
			if len(batch) == maxItems {
				sendToChannel(batches, batch)
				batch = make([]T, 0)
			}
		}
		sendToChannel(batches, batch)
	}()

	return batches
}
