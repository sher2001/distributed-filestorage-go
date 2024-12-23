package main

import (
	"fmt"

	"github.com/sher2001/go-distributed-filestorage/p2p"
)

func main() {
	fmt.Println("This is a distributed file storage")
	tr := p2p.NewTCPTransport(":8081")

	tr.ListenAndAccept()

	select {}
}
