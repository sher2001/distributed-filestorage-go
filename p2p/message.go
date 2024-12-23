package p2p

// Message holds any arbitrary data that is bieng sent over
// each transport between each nodes in a network
type Message struct {
	Payload []byte
}
