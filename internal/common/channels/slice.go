package channels

func Stream[TSrc, TDest any](src <-chan TSrc, f func(stream chan<- TDest), size int) <-chan TDest {
	res := make(chan TDest, size)
	go func() {
		defer close(res)
		f(res)
	}()
	return res
}
