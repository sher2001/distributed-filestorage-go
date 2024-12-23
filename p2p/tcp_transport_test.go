package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testTCPTransport(t *testing.B) {
	addr := ":8000"
	tr := NewTCPTransport(addr)
	assert.Equal(t, tr.listenAddr, addr)

	// Server
	assert.Nil(t, tr.ListenAndAccept())
}
