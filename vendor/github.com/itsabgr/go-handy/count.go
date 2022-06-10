package handy

func Count[T any](inputs []T, fn func(T) int) int {
	t := 0
	for _, i := range inputs {
		t += fn(i)
	}
	return t
}
