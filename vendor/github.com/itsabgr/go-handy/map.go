package handy

func Map[From any, To any](inputs []From, fn func(f From) To) (outputs []To) {
	o := make([]To, 0, len(inputs))
	for _, i := range inputs {
		o = append(o, fn(i))
	}
	return o
}
