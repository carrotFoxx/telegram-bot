package tcpServer

import (
	"encoding/hex"
	"log"
	"net"
)

type TcpServer struct {
	Addr string
}

type MsgPakage struct {
	From      string
	Conn      net.Conn
	Txt       string
	Connected bool
}

func (s TcpServer) Open(returnAnswers chan MsgPakage) error {

	addr := s.Addr
	if addr == "" {
		addr = ":8080"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
			continue
		}
		go s.handleConnection(conn, returnAnswers) //TODO: Implement me

	}
}

func (s TcpServer) handleConnection(conn net.Conn, anwser chan MsgPakage) {
	const bufsize = 1024
	buf := make([]byte, bufsize)

	for {
		reqLen, err := conn.Read(buf)
		if err != nil {
			var m MsgPakage
			m.From = conn.RemoteAddr().String()
			m.Conn = conn
			m.Connected = false
			anwser <- m
			return
		}
		_ = hex.EncodeToString(buf[0:reqLen])
		var m MsgPakage
		m.From = conn.RemoteAddr().String()
		m.Conn = conn
		m.Txt = hex.EncodeToString(buf[0:reqLen])
		m.Connected = true
		anwser <- m
	}

}
