package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cryptonetworking/cns/pkg/cns"
	"github.com/cryptonetworking/cryptography"
	"github.com/cryptonetworking/cryptography/pkg/ed25519"
	"github.com/itsabgr/go-ctx"
	"github.com/itsabgr/go-handy"
	"log"
	"net/netip"
	"os"
	"os/signal"
	"time"
)

var flagTTL = flag.Duration("ttl", time.Hour*24*30, "default record ttl")
var flagPing = flag.Duration("ping", time.Minute, "ping ttl")
var flagDir = flag.String("dir", "", "data directory")
var flagAddr = flag.String("addr", "0.0.0.0:9999", "listening address")
var aCtx, cancel = ctx.WithCancel(context.Background())

func init() {
	flag.Parse()
}
func init() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Kill, os.Interrupt)
		cancel(fmt.Errorf("SIGNAL %q", (<-c).String()))
	}()
}
func main() {
	defer handy.Catch(func(recovered error) {
		log.Fatalln(recovered)
		os.Exit(1)
	})
	defer cancel(nil)
	handy.Throw(cns.Listen(aCtx,
		netip.MustParseAddrPort(*flagAddr),
		*flagDir,
		cryptography.Gen(ed25519.Algo),
		*flagTTL,
		*flagPing,
	))
}
