package main

import (
	"log"
	"time"

	"github.com/sher2001/go-distributed-filestorage/p2p"
)

// TODO
// func OnPeer(peer p2p.Peer) error {
// 	peer.Close()
// 	fmt.Println("doing some logic with peer outside of tcp transport")
// 	return nil
// }

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":8081",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO OnPeer
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "8081_network",
		PathTransFormFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}
	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(5 * time.Second)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
