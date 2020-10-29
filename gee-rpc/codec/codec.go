package codec

import "io"

// 消息头结构
type Header struct {
	ServiceMethod string // format "Service.Method"	方法名
	Seq           uint64 // sequence number chosen by client 客户端选择的序列号
	Error         string // Error Message
}

// 消息编码接口
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

// 可根据不同的type得到不同的Codec构造函数
var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
