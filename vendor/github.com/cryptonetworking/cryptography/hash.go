package cryptography

import (
	"crypto/sha512"
	"github.com/spaolacci/murmur3"
)

func Hash(b []byte) []byte {
	hash := sha512.Sum512(b)
	return hash[:]
}

func Hash64(b []byte) uint64 {
	return murmur3.Sum64(b)
}
