package main

import (
	"fmt"
	"log"

	"github.com/sher2001/go-distributed-filestorage/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("doing some logic with peer outside of tcp transport")
	return nil
}

func main() {
	addr := ":8081"
	opts := p2p.TCPTransportOpts{
		ListenAddr:    addr,
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(opts)

	// Temporary To test
	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("message : %+v\n", msg)
		}
	}()

	fmt.Printf("This is a distributed file storage, started at %s\n", addr)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
