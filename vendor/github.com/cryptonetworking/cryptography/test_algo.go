package cryptography

import "testing"

func TestAlgo(t *testing.T, algo SigningAlgo[any, any, any]) {
	Algo := algo.Algo()
	Register(algo)
	sk := New(Algo)
	if sk.Algo() != Algo {
		t.FailNow()
	}
	msg := Rand(int(RandN(0, 512)))
	sig, err := DecodeSig(sk.Sign(msg).Encode())
	if err != nil {
		t.Fatal(err)
	}
	if sig.Algo() != Algo {
		t.FailNow()
	}
	sk, err = DecodeSK(sk.UnsafeEncode())
	if err != nil {
		t.Fatal(err)
	}
	if sk.Algo() != Algo {
		t.FailNow()
	}
	pk, err := DecodePK(sk.PK().Encode())
	if err != nil {
		t.Fatal(err)
	}
	if pk.Algo() != Algo {
		t.FailNow()
	}
	err = sig.Verify(pk, msg)
	if err != nil {
		t.Fatal(err)
	}
	err = pk.Verify(sig, msg)
	if err != nil {
		t.Fatal(err)
	}
	return
}
