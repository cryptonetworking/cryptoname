package handy

//go:noinline
func Recover[T any]() T {
	return Cast[T](recover())
}
