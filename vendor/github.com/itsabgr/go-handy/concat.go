package handy

func Concat(slices ...[]byte) []byte {
	total := 0
	for _, s := range slices {
		total += len(s)
	}
	b := make([]byte, 0, total)
	for _, s := range slices {
		b = append(b, s...)
	}
	return b
}
