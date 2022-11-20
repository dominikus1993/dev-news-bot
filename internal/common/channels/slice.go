package channels

func ToSlice[T any](s <-chan T) []T {
	res := make([]T, 0)
	for v := range s {
		res = append(res, v)
	}
	return res
}

func FromSlice[T any](s ...T) chan T {
	res := make(chan T)
	go func() {
		for _, v := range s {
			res <- v
		}
		close(res)
	}()
	return res
}
