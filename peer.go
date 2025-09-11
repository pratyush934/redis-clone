package main

import (
	"fmt"
	"log/slog"
	"net"
)

type Peer struct {
	conn net.Conn
	msg  chan []byte
}

func NewPeer(conn net.Conn, msgCh chan []byte) *Peer {
	return &Peer{
		conn: conn,
		msg:  msgCh,
	}
}

func (p *Peer) readLoop() error {

	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)

		if err != nil {
			slog.Error("peer read error, ", "err ", err)
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msg <- msgBuf
		fmt.Println(string(buf[:n]))
	}

}
