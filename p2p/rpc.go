package p2p

import "net"

// RPC holds any arbitrary data that is bieng sent over
// each transport between each nodes in a network
type RPC struct {
	From    net.Addr
	Payload []byte
}
