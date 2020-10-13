package server

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"

	"github.com/vicxu416/gocat/server/protocol"
)

func newTransPool(s *Server) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return &Transport{
				Server: s,
			}
		},
	}
}

type Transport struct {
	conn       net.Conn
	reader     *bufio.Reader
	responseCh chan *protocol.Response
	ctx        context.Context
	ctxCancel  context.CancelFunc
	*Server
}

func (tran *Transport) Do() error {
	go tran.response()
	for {
		req, err := tran.readReq()
		if err != nil {
			tran.ctxCancel()
			return err
		}
		go tran.handleReq(req)
	}
}

func (tran *Transport) handleReq(req *protocol.Request) {
	resp, err := tran.handler.HandleCmd(req)
	if err != nil {
		return
	}
	tran.responseCh <- resp
}

func (tran *Transport) setConn(conn net.Conn) {
	tran.conn = conn
	tran.reader = bufio.NewReader(conn)
}

func (tran *Transport) response() {
	for {
		select {
		case resp := <-tran.responseCh:
			tran.conn.Write(resp.ToBytes())
		case <-tran.ctx.Done():
			return
		}
	}
}

func (tran *Transport) readReq() (*protocol.Request, error) {
	tye, err := tran.reader.ReadByte()
	if err != nil && err != io.EOF {
		return nil, err
	}

	return protocol.ParseRequest(tye, tran.reader)
}
