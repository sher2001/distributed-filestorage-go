package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/sher2001/go-distributed-filestorage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransFormFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransFormFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

type Message struct {
	From    string
	payload any
}

type DataMessage struct {
	key  string
	data []byte
}

func (s *FileServer) broadcast(msg *Message) error {
	// as Peer implements Conn, it also implements Writer and Reader as well
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	fmt.Printf("while broadcasting msg: %+v\n", msg.payload)
	return gob.NewEncoder(mw).Encode(msg)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. store data to disk

	// Reader will be empty after compleatly red so need to use this.
	buff := new(bytes.Buffer)
	tee := io.TeeReader(r, buff)

	if err := s.store.Write(key, tee); err != nil {
		return err
	}

	// 2. broadcast5 this file to all known peers in the Network.

	if _, err := io.Copy(buff, r); err != nil {
		return err
	}
	dataMessage := DataMessage{
		key:  key,
		data: buff.Bytes(),
	}
	fmt.Printf("bef0re broadcasting msg: %+v\n", dataMessage)
	return s.broadcast(&Message{
		From:    "todo",
		payload: dataMessage,
	})
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to quit action")
		s.Transport.Close()
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			fmt.Println("recieved msg")
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("decoded msg: %+v\n", msg)
			if err := s.handlePayload(&msg); err != nil {
				log.Fatal(err)
			}
			fmt.Println("handled msg")
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handlePayload(p *Message) error {
	switch v := p.payload.(type) {
	case DataMessage:
		fmt.Printf("recieved data %+v\n", v)
	}
	return nil
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote: %s", p.RemoteAddr())

	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			fmt.Println("attempting to connect with: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}

	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()
	s.loop()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}
