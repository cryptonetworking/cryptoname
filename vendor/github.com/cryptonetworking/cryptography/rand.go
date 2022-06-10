package cryptography

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

func Rand(len int) []byte {
	b := make([]byte, len)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	return b
}

func RandN(min, max uint64) uint64 {
	n := binary.LittleEndian.Uint64(Rand(8))
	return (n % max) + min
}
