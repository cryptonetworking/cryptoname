package main

import (
	"context"
	"encoding/base32"
	"fmt"
	"github.com/cryptonetworking/crpyotname/pkg/datagram"
	"github.com/cryptonetworking/cryptography"
	"github.com/cryptonetworking/cryptography/pkg/ed25519"
	"github.com/itsabgr/go-handy"
	"github.com/posener/cmd"
	"log"
	"net/netip"
	"os"
	"os/signal"
	"time"
)

var args = cmd.New(cmd.OptName("cryptoname"), cmd.OptDetails("crypto-safe name service client"))
var globalCtx, cancel = context.WithCancel(context.Background())
var flagDebug = args.Bool("debug", false, "is debug")
var flagLog = args.String("log", "", "log output (empty means std)")

//
var subCmdNode = args.SubCommand("node", "initial a cryptoname node")
var flagNodeTTL = subCmdNode.Duration("ttl", handy.Month, "default record TTL")
var flagNodePingTTL = subCmdNode.Duration("ping-ttl", time.Minute, "default ping TTL")
var flagNodeDir = subCmdNode.String("dir", "", "data directory (empty means in-memory)")
var flagNodeAddr = subCmdNode.String("addr", "0.0.0.0:9999", "listening ip addr and port")
var flagNodePeer = subCmdNode.String("peer", "", "bootstrap peer")
var flagNodeSK = subCmdNode.String("sk", "", "node sk")

//
var subCmdSet = args.SubCommand("set", "set a record")
var flagSetSK = subCmdSet.String("sk", "", "secret key in base32")
var flagSetAddr = subCmdSet.String("addr", "", "record address value")
var flagSetRev = subCmdSet.Int("rev", 0, "record revision value")
var flagSetNode = subCmdSet.String("node", "", "target node address")
var flagSetTO = subCmdSet.Duration("to", time.Minute, "timeout")

//
var subCmdGet = args.SubCommand("get", "get a record")
var flagGetPK = subCmdGet.String("pk", "", "public key in base32")
var flagGetNode = subCmdGet.String("node", "", "target node address")
var flagGetTO = subCmdGet.Duration("to", time.Minute, "timeout")

//
var subCmdDerive = args.SubCommand("derive", "derive a public key from secret key")
var flagDeriveSK = subCmdDerive.String("sk", "", "target secret key")

var subCmdGen = args.SubCommand("gen", "generate a secret key")
var flagGenAlgo = subCmdGen.String("algo", ed25519.Algo, "secret key algorithm")

func main() {
	defer handy.Catch(func(recovered any) {
		if *flagDebug {
			panic(recovered)
		}
		if subCmdNode.Parsed() {
			log.Fatal(recovered)
		} else {
			fmt.Println(recovered)
			os.Exit(1)
		}
	})
	handy.Throw(args.ParseArgs(os.Args...))
	if *flagLog != "" {
		out, err := os.OpenFile(*flagLog, os.O_CREATE|os.O_APPEND, 0666)
		handy.Throw(err)
		log.SetOutput(out)
	}
	switch {
	case subCmdGen.Parsed():
		fmt.Println(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(cryptography.Gen(*flagGenAlgo).UnsafeEncode()))
	case subCmdGet.Parsed():
		globalCtx, cancel = context.WithTimeout(globalCtx, *flagGetTO)
		defer cancel()
		pk, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(*flagGetPK)
		handy.Throw(err)
		decoded, err := cryptography.DecodePK(pk)
		handy.Throw(err)
		record, err := crpyotname.Search(globalCtx, decoded, netip.MustParseAddrPort(*flagGetNode))
		handy.Throw(err)
		log.Println("Addr", record.Addr)
		log.Println("Rec", record.Revision)
	case subCmdSet.Parsed():
		globalCtx, cancel = context.WithTimeout(globalCtx, *flagSetTO)
		defer cancel()
		sk, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(*flagSetSK)
		handy.Throw(err)
		decoded, err := cryptography.DecodeSK(sk)
		handy.Throw(err)
		handy.Throw(crpyotname.Update(globalCtx, 1, decoded, uint64(*flagSetRev), netip.MustParseAddrPort(*flagSetAddr), netip.MustParseAddrPort(*flagSetNode)))
	case subCmdNode.Parsed():
		var sk *cryptography.SK
		if *flagNodeSK == "" {
			sk = cryptography.Gen(ed25519.Algo)
		} else {
			decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(*flagNodeSK)
			handy.Throw(err)
			sk, err = cryptography.DecodeSK(decoded)
			handy.Throw(err)
		}
		var peers []netip.AddrPort
		if *flagNodePeer != "" {
			peers = append(peers, netip.MustParseAddrPort(*flagNodePeer))
		}
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Kill, os.Interrupt)
			log.Println("SIGNAL", <-sigChan)
			cancel()
		}()
		dgram, err := datagram.Create(netip.MustParseAddrPort(*flagNodeAddr))
		handy.Throw(err)
		defer handy.Just(dgram.Close)
		st, err := crpyotname.OpenStorage(*flagNodeDir)
		handy.Throw(err)
		defer handy.Just(st.Close)
		node := crpyotname.New(dgram, st, sk)
		node.DefaultTTL = *flagNodeTTL
		node.PingTTL = *flagNodePingTTL
		if *flagNodePeer != "" {
			handy.Throw(node.AddPeer(netip.MustParseAddrPort(*flagNodePeer), true))
		}
		handy.Throw(node.Start(globalCtx))
	case subCmdDerive.Parsed():
		sk, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(*flagDeriveSK)
		handy.Throw(err)
		decoded, err := cryptography.DecodeSK(sk)
		handy.Throw(err)
		fmt.Println(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(decoded.PK().Encode()))
	default:
		//args.PrintDefaults()
	}
}
