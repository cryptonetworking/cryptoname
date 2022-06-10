package handy

//Box call fn and recover a panic and returns recovered value
//go:noinline
func Box(fn func()) (recovered interface{}) {
	defer func() {
		recovered = recover()
	}()
	fn()
	return
}
