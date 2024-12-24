package p2p

// Peer is an interface that represents remote nodes
type Peer interface {
	Close() error
}

// Transport is an interface that represents the communication
// layer between nodes in network (TCP, UDP, WebSockets etc.,)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
