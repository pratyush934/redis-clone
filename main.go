package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
)

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
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultAddress
	}
	return &Server{
		config: cfg,
		peers:  make(map[*Peer]bool),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.config.ListenAddress)

	fmt.Print("Hey I am Running")

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

		case <-s.quitChan:
			return
		case peer := <-s.addPeerChan:
			s.peers[peer] = true

		case peer := <-s.delPeerChan:
			delete(s.peers, peer)

		}
	}
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
	peer := NewPeer(conn)
	s.addPeerChan <- peer

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
