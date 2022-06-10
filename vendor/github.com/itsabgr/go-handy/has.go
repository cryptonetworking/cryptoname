package handy

func Contains[T any](inputs []T, fn func(T) bool) bool {
	for _, i := range inputs {
		if fn(i) {
			return true
		}
	}
	return false
}
