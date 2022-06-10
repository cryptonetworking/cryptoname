package handy

//Catch call fn with recovered value is it's not nil
func Catch[T any](fn func(recovered T)) {
	recovered := recover()
	if recovered != nil {
		fn(Cast[T](recovered))
	}
}
