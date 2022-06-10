package ed25519

import (
	lib "crypto/ed25519"
	"crypto/rand"
	"errors"
	"github.com/cryptonetworking/cryptography"
	"github.com/itsabgr/go-handy"
)

type ed25519 struct{}

const Algo = "ed25519"

func (ed25519) Algo() string {
	return Algo
}

func (ed25519) DecodeSK(b []byte) (any, error) {
	if len(b) != lib.PrivateKeySize {
		return nil, errors.New("invalid bytes length")
	}
	return lib.PrivateKey(b), nil
}

func (ed25519) DecodePK(b []byte) (any, error) {
	if len(b) != lib.PublicKeySize {
		return nil, errors.New("invalid bytes length")
	}

	return lib.PublicKey(b), nil
}

func (ed25519) DecodeSig(b []byte) (any, error) {
	if len(b) != lib.SignatureSize {
		return nil, errors.New("invalid bytes length")
	}
	return b, nil
}

func (ed25519) EncodeSK(sk any) []byte {
	return sk.(lib.PrivateKey)
}

func (ed25519) EncodePK(pk any) []byte {
	return pk.(lib.PublicKey)
}

func (ed25519) EncodeSig(sig any) []byte {
	return sig.([]byte)
}

func (ed25519) Sign(sk any, bytes []byte) any {
	return lib.Sign(sk.(lib.PrivateKey), bytes)
}

func (ed25519) Derive(sk any) any {
	return lib.PublicKey(sk.(lib.PrivateKey)[lib.PrivateKeySize-lib.PublicKeySize:])
}

func (ed25519) New() any {
	_, sk, err := lib.GenerateKey(rand.Reader)
	handy.Throw(err)
	return sk
}

func (ed25519) Verify(sig any, pk any, msg []byte) error {
	if lib.Verify(pk.(lib.PublicKey), msg, sig.([]byte)) {
		return nil
	}
	return errors.New("invalid signature")
}

func init() {
	cryptography.Register(ed25519{})
}
