package p2p

import "net"

// Peer is an interface that represents remote nodes
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport is an interface that represents the communication
// layer between nodes in network (TCP, UDP, WebSockets etc.,)
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
