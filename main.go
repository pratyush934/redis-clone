package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
)

// comment added 

const defaultAddress = ":5001"

type Config struct {
	ListenAddress string
}

type Server struct {
	config      Config
	peers       map[*Peer]bool
	ln          net.Listener
	addPeerChan chan *Peer
	delPeerChan chan *Peer
	quitChan    chan struct{}
	msg         chan []byte
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultAddress
	}
	return &Server{
		config:      cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		delPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		msg:         make(chan []byte),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.config.ListenAddress)

	fmt.Println("Hey I am Running")
	slog.Info("I am Running and you should be glad\n")

	if err != nil {
		return err
	}

	s.ln = ln

	go s.loop()

	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {

		case rawMs := <-s.msg:
			if err := s.handleRawMessage(rawMs); err != nil {
				slog.Error("issue persist in the s.handleRawMessage")
			}
			fmt.Println("rwaMsg looks like ", rawMs, "the string version is", string(rawMs))
		case <-s.quitChan:
			return
		case peer := <-s.addPeerChan:
			s.peers[peer] = true

		case peer := <-s.delPeerChan:
			delete(s.peers, peer)

		}
	}
}

func (s *Server) handleRawMessage(rawMsg []byte) error {

	fmt.Println(string(rawMsg))

	return nil
}

func (s *Server) acceptLoop() error {

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	peer := NewPeer(conn, s.msg)
	s.addPeerChan <- peer
	fmt.Println("Hello peer", peer)
	slog.Info("I am HandleConnection and I am Working", "addr", conn)

	if err := peer.readLoop(); err != nil {
		slog.Error("error is there in peer.ReadLoop", "err", err, "remoteAddr : ", conn.RemoteAddr())
	}

}

func main() {

	listenAddr := flag.String("listenAddr : ",
		defaultAddress, "listening address of the redis clone ")

	server := NewServer(Config{
		ListenAddress: *listenAddr,
	})
	log.Fatal(server.Start())
}
