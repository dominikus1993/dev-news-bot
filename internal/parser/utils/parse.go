package utils

import (
	"context"
)

func Parse[TStream any](ctx context.Context, action func(ctx context.Context, stream chan<- TStream)) <-chan TStream {
	result := make(chan TStream, 20)
	go func(context context.Context, channel chan<- TStream) {
		defer close(channel)
		action(context, channel)
	}(ctx, result)
	return result
}
