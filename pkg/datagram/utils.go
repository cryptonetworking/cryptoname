package datagram

import (
	"github.com/itsabgr/go-handy"
	"net/netip"
)

func MustGetFreeAddrPort() netip.AddrPort {
	dgram, err := Create(netip.MustParseAddrPort("0.0.0.0:0"))
	handy.Throw(err)
	defer handy.Just(dgram.Close)
	return dgram.Addr()
}
