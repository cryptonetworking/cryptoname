package cryptography

import (
	"bytes"
	"fmt"
	"github.com/itsabgr/go-handy"
	"github.com/vmihailenco/msgpack/v5"
)

type SigningAlgo[SK, PK, Sig any] interface {
	Algo() string
	DecodeSK([]byte) (SK, error)
	DecodePK([]byte) (PK, error)
	DecodeSig([]byte) (Sig, error)
	EncodeSK(SK) []byte
	EncodePK(PK) []byte
	EncodeSig(Sig) []byte
	Sign(SK, []byte) Sig
	Derive(sk SK) PK
	New() SK
	Verify(sig Sig, pk PK, msg []byte) error
}

var reg = make(map[string]SigningAlgo[any, any, any])

type SK struct {
	algo SigningAlgo[any, any, any]
	sk   any
}

type PK struct {
	algo SigningAlgo[any, any, any]
	pk   any
}
type Sig struct {
	algo SigningAlgo[any, any, any]
	sig  any
}

func (sig *Sig) Unwrap() any {
	return sig.sig
}
func (pk *PK) Unwrap() any {
	return pk.pk
}
func (sk *SK) Unwrap() any {
	return sk.sk
}
func Register(algo SigningAlgo[any, any, any]) {
	reg[algo.Algo()] = algo
}
func Gen(name string) *SK {
	algo := reg[name]
	sk := algo.New()
	return &SK{algo, sk}
}
func (sk *SK) Sign(msg []byte) *Sig {
	algo := sk.algo
	sig := algo.Sign(sk.sk, msg)
	return &Sig{algo, sig}
}

func New(algo string) *SK {
	algorithm := reg[algo]
	sk := algorithm.New()
	return &SK{algorithm, sk}
}

type encoded struct {
	Kind  string
	Algo  string
	Bytes []byte
}

func (sk *SK) UnsafeEncode() []byte {
	algo := sk.algo
	name := algo.Algo()
	b := algo.EncodeSK(sk.sk)
	buff := bytes.NewBuffer(nil)
	handy.Throw(msgpack.NewEncoder(buff).Encode(encoded{"SK", name, b}))
	return buff.Bytes()
}
func (sig *Sig) Encode() []byte {
	algo := sig.algo
	name := algo.Algo()
	b := algo.EncodeSig(sig.sig)
	buff := bytes.NewBuffer(nil)
	handy.Throw(msgpack.NewEncoder(buff).Encode(encoded{"Sig", name, b}))
	return buff.Bytes()
}
func (pk *PK) Encode() []byte {
	algo := pk.algo
	name := algo.Algo()
	b := algo.EncodePK(pk.pk)
	buff := bytes.NewBuffer(nil)
	handy.Throw(msgpack.NewEncoder(buff).Encode(encoded{"PK", name, b}))
	return buff.Bytes()
}

func DecodeSK(b []byte) (*SK, error) {
	enc := new(encoded)
	buff := bytes.NewReader(b)
	handy.Throw(msgpack.NewDecoder(buff).Decode(enc))
	algo, found := reg[enc.Algo]
	if !found {
		return nil, fmt.Errorf("unsupported algorithm %q", enc.Algo)
	}
	sk, err := algo.DecodeSK(enc.Bytes)
	if err != nil {
		return nil, err
	}
	return &SK{algo, sk}, nil
}
func DecodePK(b []byte) (*PK, error) {
	enc := new(encoded)
	buff := bytes.NewReader(b)
	handy.Throw(msgpack.NewDecoder(buff).Decode(enc))
	algo, found := reg[enc.Algo]
	if !found {
		return nil, fmt.Errorf("unsupported algorithm %q", enc.Algo)
	}
	pk, err := algo.DecodePK(enc.Bytes)
	if err != nil {
		return nil, err
	}
	return &PK{algo, pk}, nil
}
func DecodeSig(b []byte) (*Sig, error) {
	enc := new(encoded)
	buff := bytes.NewReader(b)
	handy.Throw(msgpack.NewDecoder(buff).Decode(enc))
	algo, found := reg[enc.Algo]
	if !found {
		return nil, fmt.Errorf("unsupported algorithm %q", enc.Algo)
	}
	sig, err := algo.DecodeSig(enc.Bytes)
	if err != nil {
		return nil, err
	}
	return &Sig{algo, sig}, nil
}
func (sk *SK) PK() *PK {
	algo := sk.algo
	pk := algo.Derive(sk.sk)
	return &PK{algo, pk}
}
func (pk *PK) Verify(sig *Sig, msg []byte) error {
	algo := pk.algo
	return algo.Verify(sig.sig, pk.pk, msg)
}
func (sig *Sig) Verify(pk *PK, msg []byte) error {
	algo := sig.algo
	return algo.Verify(sig.sig, pk.pk, msg)
}

func (sig *Sig) Algo() string {
	return sig.algo.Algo()
}

func (pk *PK) Algo() string {
	return pk.algo.Algo()
}

func (sk *SK) Algo() string {
	return sk.algo.Algo()
}
