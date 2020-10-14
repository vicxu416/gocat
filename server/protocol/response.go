package protocol

import (
	"bufio"
	"bytes"
	"io"

	"github.com/vicxu416/gocat/server/util"
)

type Status int8

const (
	StastusOK Status = iota + 1
	StastusRedirect
	StastusNotFound
	StastusOperationError
)

var statusMap = map[byte]Status{
	'1': StastusOK,
	'2': StastusRedirect,
	'3': StastusNotFound,
	'4': StastusOperationError,
}

func readData(lenByte []byte, reader *bufio.Reader) ([]byte, error) {
	len, err := util.BytesToInt(lenByte)
	if err != nil {
		return nil, err
	}
	data := make([]byte, len)
	if _, err = io.ReadFull(reader, data); err != nil {
		return nil, err
	}
	return data, nil
}

func readHeader(reader *bufio.Reader) ([][]byte, error) {
	header, err := reader.ReadBytes(' ')
	if err != nil {
		return nil, err
	}
	header = bytes.TrimSpace(header)
	headers := bytes.Split(header, []byte(","))
	if bytes.Equal(headers[len(headers)-1], []byte{}) {
		headers = headers[:len(headers)-1]
	}

	return headers, nil
}

type Response struct {
	status Status
	key    []byte
	val    []byte
	msg    string
	errMsg string
}

func (resp *Response) Status() Status {
	return resp.status
}

func (resp *Response) Unmarshal(reader *bufio.Reader) error {
	stat, err := reader.ReadByte()
	if err != nil {
		return err
	}
	resp.status = statusMap[stat]
	headers, err := readHeader(reader)
	if err != nil {
		return err
	}
	for _, header := range headers {
		payloadType := header[0]
		switch payloadType {
		case 'k':
			if resp.key, err = readData(header[1:], reader); err != nil {
				return err
			}
		case 'v':
			if resp.val, err = readData(header[1:], reader); err != nil {
				return err
			}
		case 'm':
			msg, err := readData(header[1:], reader)
			if err != nil {
				return err
			}
			resp.msg = util.BytesToString(msg)
		case 'e':
			errMsg, err := readData(header[1:], reader)
			if err != nil {
				return err
			}
			resp.errMsg = util.BytesToString(errMsg)
		default:
		}
	}
	return nil
}

func (resp *Response) Marshal() []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 1))
	stat := util.IntToBytes((int(resp.status)))
	stat = bytes.TrimSpace(stat)
	_, _ = buffer.Write(stat)
	resp.writeHeader(buffer)
	resp.writeData(buffer)
	return buffer.Bytes()
}

func (resp *Response) SeStatus(status Status) {
	resp.status = status
}

func (resp *Response) writeData(buffer *bytes.Buffer) {
	if len(resp.key) > 0 {
		_, _ = buffer.Write(resp.key)
	}

	if len(resp.val) > 0 {
		_, _ = buffer.Write(resp.val)
	}

	if len(resp.msg) > 0 {
		_, _ = buffer.WriteString(resp.msg)
	}

	if len(resp.errMsg) > 0 {
		_, _ = buffer.WriteString(resp.errMsg)
	}
}

func (resp *Response) writeHeader(buffer *bytes.Buffer) {
	if len(resp.key) > 0 {
		buffer.WriteByte('k')
		buffer.Write(util.IntToBytes(len(resp.key)))
		buffer.WriteByte(',')
	}

	if len(resp.val) > 0 {
		buffer.WriteByte('v')
		buffer.Write(util.IntToBytes(len(resp.val)))
		buffer.WriteByte(',')
	}

	if len(resp.msg) > 0 {
		buffer.WriteByte('m')
		buffer.Write(util.IntToBytes(len(resp.msg)))
		buffer.WriteByte(',')
	}

	if len(resp.errMsg) > 0 {
		buffer.WriteByte('e')
		buffer.Write(util.IntToBytes(len(resp.errMsg)))
		buffer.WriteByte(',')
	}
	buffer.WriteByte(' ')
}

func (resp *Response) SetKey(key []byte) {
	resp.key = key
}

func (resp *Response) SetVal(val []byte) {
	resp.val = val
}

func (resp *Response) GetKey() []byte {
	return resp.key
}

func (resp *Response) GetKeyStr() string {
	return util.BytesToString(resp.key)
}

func (resp *Response) GetVal() []byte {
	return resp.val
}

func (resp *Response) GetValStr() string {
	return util.BytesToString(resp.val)
}
