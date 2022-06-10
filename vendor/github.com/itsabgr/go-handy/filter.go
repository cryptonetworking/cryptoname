package handy

func Filter[T any](inputs []T, fn func(T) bool) []T {
	o := make([]T, 0, len(inputs))
	for _, i := range inputs {
		if fn(i) {
			o = append(o, i)
		}
	}
	return o
}
