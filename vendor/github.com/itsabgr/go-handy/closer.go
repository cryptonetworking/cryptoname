package handy

func Just(closable func() error) {
	if closable != nil {
		_ = closable()
	}
}
