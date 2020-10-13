package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequest(t *testing.T) {
	testcases := []struct {
		name     string
		key, val string
		req      []byte
		tye      ReqType
	}{
		{
			name: "SET", key: "test", val: "test",
			req: append([]byte{'S', '4', ' ', '4', ' '}, []byte("testtest")...), tye: SET,
		},
		{
			name: "GET", key: "test", val: "",
			req: append([]byte{'G', '4', ' '}, []byte("test")...), tye: GET,
		},
		{
			name: "DEL", key: "test", val: "",
			req: append([]byte{'D', '4', ' '}, []byte("test")...), tye: DEL,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewReader(testcase.req[1:]))
			request, err := ParseRequest(testcase.req[0], reader)
			assert.NoError(t, err)
			assert.Equal(t, request.Typ, testcase.tye)
			assert.Equal(t, request.GetKeyStr(), testcase.key)
			assert.Equal(t, request.GetValStr(), testcase.val)
		})
	}

}
