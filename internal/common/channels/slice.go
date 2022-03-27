package channels

func ToSlice[T any](s <-chan T) []T {
	res := make([]T, 0)
	for v := range s {
		res = append(res, v)
	}
	return res
}
