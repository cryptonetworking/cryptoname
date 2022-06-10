package handy

//Throw panics if any of values is not nil
func Throw(values ...interface{}) {
	for _, error := range values {
		if error != nil {
			panic(error)
		}
	}
}
