package handy

//Clone makes a copy of bytes
func Clone(src []byte) []byte {
	clone := make([]byte, len(src))
	Assert(copy(clone, src) == len(src))
	return clone
}
