package server

import (
	"bufio"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vicxu416/gocat/server/protocol"
)

var (
	addr = "127.0.0.1"
	port = "8888"
)

func TestTransport(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite
	server      *Server
	client      net.Conn
	reader      *bufio.Reader
	handler     Handler
	mockHandler *HandlerMock
	mock        bool
}

func (s *testSuite) writeRequest(req *protocol.Request) {
	_, err := s.client.Write(req.ToBytes())
	s.Require().NoError(err)
	time.Sleep(time.Millisecond)
}

func (s *testSuite) SetupSuite() {
	if s.handler == nil {
		s.mockHandler = new(HandlerMock)
		s.mock = true
		s.server = New(addr, port, s.mockHandler)
	} else {
		s.server = New(addr, port, s.handler)
	}

	go func() {
		if err := s.server.Run(); err != nil {
			log.Panicf("run server failed, err:%+v", err)
			os.Exit(1)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	client, err := net.Dial("tcp", addr+":"+port)
	s.Require().NoError(err)
	s.client = client
	s.reader = bufio.NewReader(client)
}

func (s *testSuite) TearDownSuite() {
	s.server.Shutdown()
}

func (s *testSuite) TestGetRequest() {
	req := protocol.NewTestRequest(protocol.GET, []byte("testing"), nil)
	if s.mock {
		s.mockHandler.On("Handle", req)
	}

	s.writeRequest(req)
	resp := &protocol.Response{}
	err := resp.Unmarshal(s.reader)
	s.Assert().NoError(err)
	s.Assert().Equal(resp.Status(), protocol.StastusOK)
	s.Assert().Equal(resp.GetKeyStr(), "testing")

	if s.mock {
		s.mockHandler.AssertExpectations(s.T())
		s.Assert().Equal(resp.GetValStr(), "testing value")
	}
}

func (s *testSuite) TestSetRequest() {
	req := protocol.NewTestRequest(protocol.SET, []byte("testing"), []byte("testing123"))
	if s.mock {
		s.mockHandler.On("Handle", req)
	}
	s.writeRequest(req)
	resp := &protocol.Response{}
	err := resp.Unmarshal(s.reader)
	s.Assert().NoError(err)
	s.Assert().Equal(resp.Status(), protocol.StastusOK)

	if s.mock {
		s.mockHandler.AssertExpectations(s.T())
	}
}

func (s *testSuite) TestDelRequest() {
	req := protocol.NewTestRequest(protocol.DEL, []byte("testing"), nil)
	if s.mock {
		s.mockHandler.On("Handle", req)
	}
	s.writeRequest(req)
	resp := &protocol.Response{}
	err := resp.Unmarshal(s.reader)
	s.Assert().NoError(err)
	s.Assert().Equal(resp.Status(), protocol.StastusOK)

	if s.mock {
		s.mockHandler.AssertExpectations(s.T())
	}
}
