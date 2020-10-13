package protocol

type Status int8

const (
	OK Status = iota + 1
	Redirect
	NotFound
	OperationError
)

type Response struct {
	status Status
	key    []byte
	val    []byte
	msg    string
	errMsg string
}

func (resp *Response) ToBytes() []byte {
	return nil
}

func (resp *Response) SetKey(key []byte) {
}

func (resp *Response) SetVal(key []byte) {
}
