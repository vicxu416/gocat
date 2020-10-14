package server

import (
	"github.com/stretchr/testify/mock"
	"github.com/vicxu416/gocat/server/protocol"
)

type HandlerMock struct {
	mock.Mock
}

func (mock *HandlerMock) Handle(req *protocol.Request) (*protocol.Response, error) {
	mock.Called(req)
	response := &protocol.Response{}
	response.SeStatus(protocol.StastusOK)

	if req.Typ == protocol.GET {
		response.SetKey(req.GetKey())
		response.SetVal([]byte("testing value"))
	}
	return response, nil
}
