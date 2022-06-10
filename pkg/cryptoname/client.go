package cryptoname

import (
	"bytes"
	"context"
	"github.com/cryptonetworking/cryptography"
	"github.com/cryptonetworking/cryptoname/pkg/datagram"
	"github.com/itsabgr/go-handy"
	"net/netip"
	"time"
)

func Update(ctx context.Context, kind uint8, sk *cryptography.SK, rev uint64, addr netip.AddrPort, nodes ...netip.AddrPort) error {
	conn, err := datagram.Create(netip.MustParseAddrPort("0.0.0.0:0"))
	if err != nil {
		return err
	}
	defer handy.Just(conn.Close)
	go func() {
		if rev == 0 {
			rev = uint64(time.Now().Unix())
		}
		update := &Record{
			Addr: addr,
			Rev:  rev,
			PK:   sk.PK().Encode(),
			Kind: kind,
		}
		update.Sign(sk)
		for ctx.Err() == nil {
			for _, node := range nodes {
				_ = conn.SendTo(update.Encode(), node)
			}
		}
	}()
	for ctx.Err() == nil {
		record, err := Search(ctx, sk.PK(), nodes...)
		if err != nil {
			return err
		}
		if record.Rev == rev {
			return nil
		}
	}
	return ctx.Err()
}
func Search(ctx context.Context, pk *cryptography.PK, nodes ...netip.AddrPort) (*Record, error) {
	conn, err := datagram.Create(netip.MustParseAddrPort("0.0.0.0:0"))
	if err != nil {
		return nil, err
	}
	defer handy.Just(conn.Close)
	go func() {
		query := (&Record{PK: pk.Encode()}).Encode()
		for ctx.Err() == nil {
			for _, node := range nodes {
				_ = conn.SendTo(query, node)
			}
		}
	}()
	b := make([]byte, 1024)
	record := new(Record)
	for ctx.Err() != nil {
		n, _, err := conn.RecvFrom(ctx, b)
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(record.PK, pk.Encode()) {
			continue
		}
		err = record.Decode(b[:n])
		if err != nil {
			continue
		}
		err = record.Verify()
		if err != nil {
			continue
		}
		return record, nil
	}
	return nil, ctx.Err()
}
