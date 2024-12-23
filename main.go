package main

import (
	"fmt"
	"log"

	"github.com/sher2001/go-distributed-filestorage/p2p"
)

func main() {
	addr := ":8081"
	opts := p2p.TCPTransportOpts{
		ListenAddr:    addr,
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransport(opts)

	fmt.Printf("This is a distributed file storage, started at %s\n", addr)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
