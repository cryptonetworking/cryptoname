package datagram

import (
	"context"
	"fmt"
	"github.com/itsabgr/go-handy"
	"golang.org/x/sys/unix"
	"net/netip"
	"runtime"
	"syscall"
)

type Datagram struct {
	fd int
}

func (dgram *Datagram) FD() int {
	return dgram.fd
}
func (dgram *Datagram) Addr() netip.AddrPort {
	addr, err := syscall.Getsockname(dgram.fd)
	handy.Throw(err)
	return sockaddrToAddrPort(addr)
}
func (dgram *Datagram) ReuseAddrPort() error {
	err := syscall.SetsockoptInt(dgram.fd, syscall.SOL_SOCKET, 0xf, 1)
	if err != nil {
		return err
	}
	return syscall.SetsockoptInt(dgram.fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
}

func New(fd int) *Datagram {
	return &Datagram{fd}
}

func Create(addr netip.AddrPort) (*Datagram, error) {
	var fd int
	var err error
	if addr.Addr().Is4() {
		fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
		if err != nil {
			return nil, err
		}
		err = syscall.Bind(fd, addrPortToSockaddr(addr))
		if err != nil {
			_ = syscall.Close(fd)
			return nil, err
		}
		return New(fd), nil
	} else {
		fd, err = syscall.Socket(syscall.AF_INET6, syscall.SOCK_DGRAM, 0)
		if err != nil {
			return nil, err
		}
		err = syscall.Bind(fd, addrPortToSockaddr(addr))
		if err != nil {
			_ = syscall.Close(fd)
			return nil, err
		}
		fmt.Println(addrPortToSockaddr(addr))
		err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.IPV6_V6ONLY, 0)
		if err != nil {
			_ = syscall.Close(fd)
			return nil, err
		}
		return New(fd), nil
	}

}

func sockaddrToAddrPort(sockaddr syscall.Sockaddr) netip.AddrPort {
	switch sockaddr.(type) {
	case *syscall.SockaddrInet4:
		sockaddr := sockaddr.(*syscall.SockaddrInet4)
		return netip.AddrPortFrom(netip.AddrFrom4(sockaddr.Addr), uint16(sockaddr.Port))
	case *syscall.SockaddrInet6:
		sockaddr := sockaddr.(*syscall.SockaddrInet6)
		return netip.AddrPortFrom(netip.AddrFrom16(sockaddr.Addr), uint16(sockaddr.Port))
	default:
		panic(fmt.Errorf("unsupported Sockaddr %+v", sockaddr))
	}
}
func addrPortToSockaddr(addrPort netip.AddrPort) syscall.Sockaddr {
	if addrPort.Addr().Is4In6() || addrPort.Addr().Is4() {
		return &syscall.SockaddrInet4{
			Port: int(addrPort.Port()),
			Addr: addrPort.Addr().As4(),
		}
	}
	return &syscall.SockaddrInet6{
		Port:   int(addrPort.Port()),
		Addr:   addrPort.Addr().As16(),
		ZoneId: 0, //TODO
	}
}
func (dgram *Datagram) Close() error {
	if dgram == nil {
		return nil
	}
	return syscall.Close(dgram.fd)
}

func (dgram *Datagram) RecvFrom(ctx context.Context, b []byte) (int, netip.AddrPort, error) {
	for ctx.Err() == nil {
		n, from, err := syscall.Recvfrom(dgram.fd, b, 0)
		if err != nil {
			if err == unix.EAGAIN {
				runtime.Gosched()
				continue
			}
			return 0, netip.AddrPort{}, err
		}
		addr := sockaddrToAddrPort(from)
		return n, addr, nil
	}
	return 0, netip.AddrPort{}, ctx.Err()
}

func (dgram *Datagram) SendTo(b []byte, to netip.AddrPort) error {
	return syscall.Sendto(dgram.fd, b, 0, addrPortToSockaddr(to))
}
