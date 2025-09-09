package main

import (
	"fmt"
	"net"
)

type Peer struct {
	conn net.Conn
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		conn: conn,
	}
}

func (p *Peer) readLoop() error {
	for {
		fmt.Print("read ho ho")
	}
	return nil
}
