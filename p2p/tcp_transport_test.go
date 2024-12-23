package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	addr := ":4000"
	opts := TCPTransportOpts{
		ListenAddr: addr,
	}
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddr, addr)

	// Server
	assert.Nil(t, tr.ListenAndAccept())
}
