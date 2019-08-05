package reader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// 多重读取器接口
type MultipleReader interface {
	// 获取一个可关闭读取器的实例
	Reader() io.ReadCloser
}

// 多重读取器的实现
type myMultipleReader struct {
	data []byte
}

// 创建并返回一个多重读取器的实例
func NewMultipleReader(reader io.Reader) (MultipleReader, error) {
	var data []byte
	var err error
	if reader != nil {
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("multiple reader: couldn't create a new one: %s", err)
		}
	} else {
		data = []byte{}
	}
	return &myMultipleReader{
		data:  data,
	}, nil
}

func (rr *myMultipleReader) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(rr.data))
}
