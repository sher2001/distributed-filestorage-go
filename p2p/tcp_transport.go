package p2p

import (
	"errors"
	"fmt"
	"net"
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

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

// RemoteAddr implements the Peer Interface and will return the
// remote address of its underlying connection.
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

// close implements peer intefacce.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	fmt.Printf("TCP transport listening on port %s\n", t.ListenAddr)

	return nil
}

// consume inplements the Transport interface and will
// returns read-only channel for reading the incomming
// messages recieved from another peer in tghe network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// Dial implements the Transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP Accept error: %s\n", err)
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.HandShakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP Handshake error : %s\n", err)
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read Loop
	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP Decode error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

}
