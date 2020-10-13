package server

import (
	"net"
	"sync"

	"github.com/vicxu416/gocat/server/protocol"
)

type Handler interface {
	HandleCmd(cmd *protocol.Request) (*protocol.Response, error)
}

func New(addr, port string, handler Handler, options ...Option) *Server {
	config := &Config{}

	for _, opt := range options {
		opt(config)
	}

	s := &Server{
		addr:    addr,
		port:    port,
		handler: handler,
		Config:  config,
	}
	s.transPool = newTransPool(s)
	return s
}

type Config struct {
	log Logger
}

type Server struct {
	*Config
	addr      string
	port      string
	handler   Handler
	transPool *sync.Pool
}

func (s *Server) Listen() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.addr+":"+s.port)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	trans := s.transPool.Get().(*Transport)
	trans.setConn(conn)
	trans.Do()
	s.transPool.Put(trans)
}
