package codec

import (
	"bytes"
	"encoding/gob"
	"github.com/vmihailenco/msgpack"
)

type SerializeType byte

const (
	MessagePack SerializeType = iota
	GOB
)

type Codec interface {
	Encode(value interface{}) ([]byte, error)
	Decode(data []byte, value interface{}) error
}

var codecs = map[SerializeType]Codec {
	MessagePack: &MessagePackCodec{},
	GOB:         &GobCodec{},
}

type MessagePackCodec struct {}

type GobCodec struct {}

func GetCodec(t SerializeType) Codec {
	return codecs[t]
}

func (c *MessagePackCodec) Encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (c *MessagePackCodec) Decode(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (g *GobCodec) Encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	return buf.Bytes(), err
}

func (g *GobCodec) Decode(data []byte, value interface{}) error {
	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(value)
	return err
}
