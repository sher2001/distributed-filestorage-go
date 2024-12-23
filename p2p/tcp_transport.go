package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over
// TCP established connection
type TCPPeer struct {
	// conn representgs underlying connection of a peer
	conn net.Conn

	// If we dail and retrive a conn, outbound = true
	// If we accept and retrive a conn, outbound = false (inbound)
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Accept error: %s\n", err)
		}

		fmt.Printf("new incomming connection : %+v\n", conn)
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandShakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP Handshake error : %s\n", err)
		return
	}

	// Read Loop
	rpc := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, rpc); err != nil {
			fmt.Printf("TCP Decode error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		fmt.Printf("message: %+v\n", rpc)
	}

}
