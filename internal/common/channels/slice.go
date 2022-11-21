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
		defer close(res)
		for _, v := range s {
			res <- v
		}
	}()
	return res
}

func Stream[TSrc, TDest any](src <-chan TSrc, f func(stream chan<- TDest), size int) <-chan TDest {
	res := make(chan TDest, size)
	go func() {
		defer close(res)
		f(res)
	}()
	return res
}
