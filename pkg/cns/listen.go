package cns

import (
	"context"
	"github.com/cryptonetworking/cns/pkg/datagram"
	"github.com/cryptonetworking/cryptography"
	"github.com/itsabgr/go-handy"
	"net/netip"
	"time"
)

func Listen(ctx context.Context, addr netip.AddrPort, storage string, id *cryptography.SK, DefaultTTL, PingTTL time.Duration, peers ...netip.AddrPort) error {
	dgram, err := datagram.Create(addr)
	if err != nil {
		return err
	}
	defer handy.Just(dgram.Close)
	st, err := OpenStorage(storage)
	if err != nil {
		return err
	}
	defer handy.Just(st.Close)
	node := New(dgram, st, id)
	node.id = id
	node.DefaultTTL = DefaultTTL
	node.PingTTL = PingTTL
	for _, addr := range peers {
		err := node.AddPeer(addr, true)
		if err != nil {
			return err
		}
	}
	return node.Start(ctx)
}
