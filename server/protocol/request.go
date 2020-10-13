package protocol

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"strconv"

	"github.com/vicxu416/gocat/server/util"
)

type ReqType int8

const (
	GET ReqType = iota + 1
	SET
	DEL
)

var cmdMap = map[byte]ReqType{
	'S': SET,
	'G': GET,
	'D': DEL,
}

type Request struct {
	Typ ReqType
	key []byte
	val []byte
	ctx context.Context
}

func (cmd *Request) GetKey() []byte {
	return cmd.key
}

func (cmd *Request) GetKeyStr() string {
	return util.BytesToString(cmd.key)
}

func (cmd *Request) GetVal() []byte {
	return cmd.val
}

func (cmd *Request) GetValStr() string {
	return util.BytesToString(cmd.val)
}

func ParseRequest(tye byte, reader *bufio.Reader) (*Request, error) {
	var (
		req = &Request{}
		err error
	)
	req.Typ = cmdMap[tye]
	switch req.Typ {
	case SET:
		err = readKeyVal(reader, req)
	case GET, DEL:
		err = readKey(reader, req)
	default:
	}

	if err != nil {
		return nil, err
	}
	return req, nil
}

func readKeyVal(reader *bufio.Reader, req *Request) error {
	keyLen, err := readLen(reader)
	if err != nil {
		return err
	}
	valLen, err := readLen(reader)
	if err != nil {
		return err
	}
	keyByt := make([]byte, keyLen)
	valByt := make([]byte, valLen)

	if _, err = io.ReadFull(reader, keyByt); err != nil {
		return err
	}

	if _, err = io.ReadFull(reader, valByt); err != nil {
		return err
	}

	req.key = keyByt
	req.val = valByt
	return nil
}

func readKey(reader *bufio.Reader, req *Request) error {
	keyLen, err := readLen(reader)
	if err != nil {
		return err
	}
	keyByt := make([]byte, keyLen)
	if _, err = io.ReadFull(reader, keyByt); err != nil {
		return err
	}
	req.key = keyByt
	return nil
}

func readLen(reader *bufio.Reader) (int, error) {
	lenBty, err := reader.ReadBytes(' ')
	if err != nil {
		return 0, err
	}
	lenBty = bytes.TrimSpace(lenBty)
	return strconv.Atoi(util.BytesToString(lenBty))
}
