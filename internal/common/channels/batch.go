package channels

import (
	"context"
	"sync"
)

type ErrorStream[T any] struct {
	stream chan T
	errors chan error
}

func NewErrorStream[T any](size int) *ErrorStream[T] {
	return &ErrorStream[T]{
		stream: make(chan T, size),
		errors: make(chan error, size),
	}
}

func (stream *ErrorStream[T]) Close() {
	defer close(stream.stream)
	defer close(stream.errors)
}

func (stream *ErrorStream[T]) SendError(err error) {
	stream.errors <- err
}

func (stream *ErrorStream[T]) Send(el T) {
	stream.stream <- el
}

func (stream *ErrorStream[T]) SendArr(elements []T) {
	for _, v := range elements {
		stream.stream <- v
	}
}

func (stream *ErrorStream[T]) Read() (<-chan T, <-chan error) {
	return stream.stream, stream.errors
}

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

func RunBatchInPararell[T any](ctx context.Context, values <-chan T, maxItems, size int, f func(context.Context, []T) ([]T, error)) *ErrorStream[T] {
	result := NewErrorStream[T](size)
	go func() {
		defer result.Close()
		var wg sync.WaitGroup
		for value := range Batch(ctx, values, maxItems) {
			arr := value
			wg.Go(func() {
				res, err := f(ctx, arr)
				if err != nil {
					result.SendError(err)
				} else {
					result.SendArr(res)
				}
			})
		}
		wg.Wait()
	}()
	return result
}
