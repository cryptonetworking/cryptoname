package cryptoname

import (
	"bytes"
	"errors"
	"github.com/cryptonetworking/cryptography"
	"github.com/itsabgr/go-handy"
	"net/netip"
)
import "github.com/vmihailenco/msgpack/v5"

type Record struct {
	Sig  []byte         `msgpack:"sig,omitempty"`
	PK   []byte         `msgpack:"pk,omitempty"`
	Kind uint8          `msgpack:"kind,omitempty"`
	Rev  uint64         `msgpack:"rev,omitempty"`
	Addr netip.AddrPort `msgpack:"addr,omitempty"`
}

func (r *Record) Encode() []byte {
	b := bytes.NewBuffer(nil)
	handy.Throw(msgpack.NewEncoder(b).Encode(r))
	return b.Bytes()
}
func (r *Record) Decode(b []byte) error {
	return msgpack.NewDecoder(bytes.NewReader(b)).Decode(r)
}
func (r *Record) Sign(sk *cryptography.SK) {
	r.PK = sk.PK().Encode()
	r.Sig = sk.Sign(r.digest()).Encode()
}

func (r *Record) digest() []byte {
	r2 := *r
	r2.Sig = nil
	b := bytes.NewBuffer(nil)
	handy.Throw(msgpack.NewEncoder(b).Encode(&r2))
	return b.Bytes()
}

func (r *Record) Verify() error {
	if !r.Addr.Addr().IsValid() {
		return errors.New("invalid addr")
	}
	pk, err := cryptography.DecodePK(r.PK)
	if err != nil {
		return err
	}
	sig, err := cryptography.DecodeSig(r.Sig)
	if err != nil {
		return err
	}
	return pk.Verify(sig, r.digest())
}
