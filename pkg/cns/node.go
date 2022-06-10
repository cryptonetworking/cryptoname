package cns

import (
	"context"
	"errors"
	"github.com/cryptonetworking/cns/pkg/datagram"
	"github.com/cryptonetworking/cryptography"
	"github.com/cryptonetworking/cryptography/pkg/ed25519"
	"github.com/itsabgr/go-ctx"
	"github.com/itsabgr/go-handy"
	"github.com/samber/lo"
	"net/netip"
	"runtime"
	"time"
)

type ping struct {
	deadline time.Time
	addr     netip.AddrPort
	pin      bool
}

type node struct {
	DefaultTTL time.Duration
	PingTTL    time.Duration
	id         *cryptography.SK
	conn       *datagram.Datagram
	storage    *storage
	pings      []*ping
}

func (n *node) Start(aCtx context.Context) error {
	aCtx, cancel := ctx.WithCancel(aCtx)
	defer cancel(nil)
	for range handy.N(runtime.NumCPU()) {
		go func() {
			err := n.listen(aCtx)
			if err != nil {
				cancel(err)
			}
		}()
	}
	return n.pinger(aCtx)
}
func New(conn *datagram.Datagram, storage *storage) *node {
	n := &node{conn: conn, storage: storage, id: cryptography.Gen(ed25519.Algo)}
	return n
}
func (n *node) pinger(ctx context.Context) error {
	for ctx.Err() == nil {
		_ = n.ping()
		<-time.NewTimer(time.Second).C
	}
	return ctx.Err()
}

func (n *node) AddPeer(addrPort netip.AddrPort, pin bool) error {
	toString := addrPort.String()
	if n.conn.Addr().String() == toString {
		return nil
	}
	p, _, ok := lo.FindIndexOf(n.pings, func(p *ping) bool {
		return p.addr.String() == toString
	})
	if ok {
		p.deadline = time.Now().Add(n.PingTTL)
		p.pin = pin
		return nil
	}
	n.pings = append(n.pings, &ping{
		deadline: time.Now().Add(n.PingTTL),
		addr:     addrPort,
		pin:      pin,
	})
	return nil
}

func (n *node) Peers() []netip.AddrPort {
	now := time.Now()
	n.pings = handy.Filter(n.pings, func(p *ping) bool {
		if p.pin {
			return true
		}
		return p.deadline.After(now)
	})
	return handy.Map(n.pings, func(p *ping) netip.AddrPort {
		return p.addr
	})
}
func (n *node) broadcast(record *Record, addrs ...netip.AddrPort) error {
	lastVer, err := n.storage.GetVersion(record)
	if err != nil {
		return err
	}
	packet := record.Encode()
	for _, addr := range addrs {
		ver, err := n.storage.GetVersion(record)
		if err != nil {
			return err
		}
		if ver != lastVer {
			return errors.New("record changed")
		}
		_ = n.conn.SendTo(packet, addr)
	}
	return nil
}
func (n *node) ping() error {
	record := new(Record)
	record.Revision = uint64(time.Now().Unix())
	record.Addr = n.conn.Addr()
	record.Sign(n.id)
	err := n.storage.Store(record, n.PingTTL)
	if err != nil {
		return err
	}
	return n.broadcast(record, n.Peers()...)
}
func (n *node) listen(ctx context.Context) error {
	b := make([]byte, 1024)
	record := new(Record)
	for {
		readN, from, err := n.conn.RecvFrom(ctx, b)
		if err != nil {
			return err
		}
		err = record.Decode(b[:readN])
		if err != nil {
			continue
		}
		if len(record.Sig) == 0 {
			err := n.storage.Load(record)
			if err == nil {
				_ = n.conn.SendTo(record.Encode(), from)
				continue
			}
			_ = n.broadcast(record, n.Peers()...)
			continue
		}
		err = record.Verify()
		if err != nil {
			continue
		}
		var ttl time.Duration
		if record.Kind == 0 {
			err = n.AddPeer(record.Addr, false)
			if err != nil {
				continue
			}
			ttl = n.PingTTL
		} else {
			ttl = n.DefaultTTL
		}
		err = n.storage.Store(record, ttl)
		if err != nil {
			if err == ErrOldRecord {
				err = n.storage.Load(record)
				if err == nil {
					_ = n.conn.SendTo(record.Encode(), from)
				}
			}
			continue
		}
		_ = n.broadcast(record, n.Peers()...)
	}
}
