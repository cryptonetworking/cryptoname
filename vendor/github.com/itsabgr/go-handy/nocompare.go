package handy

//NoCompare prevents struct to be compared
type NoCompare struct {
	NotComparable [0]func() `json:"-"`
}
