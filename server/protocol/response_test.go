package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseToBytes(t *testing.T) {
	resp := &Response{
		status: StastusOK,
		key:    []byte("testing"),
		val:    []byte("testing12314312341235234safsdf234sdf"),
		msg:    "hello",
		errMsg: "test error msg",
	}

	respByt := resp.Marshal()
	reader := bufio.NewReader(bytes.NewReader(respByt))
	resp2 := &Response{}
	err := resp2.Unmarshal(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp2.status, resp.status)
	assert.Equal(t, resp2.key, resp.key)
	assert.Equal(t, resp2.val, resp.val)
	assert.Equal(t, resp2.msg, resp.msg)
	assert.Equal(t, resp2.errMsg, resp.errMsg)
}
