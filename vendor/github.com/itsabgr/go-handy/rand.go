package handy

import (
	"crypto/rand"
	"io"
)

func Rand(length int) []byte {
	b := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, b)
	Throw(err)
	return b
}
