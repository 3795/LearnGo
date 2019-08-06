package buffer

import (
	"Project/LearnGo/Learn/webcrawler/errors"
	er "errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// ErrClosedBufferPool 是表示缓冲池已关闭的错误的变量。
var ErrClosedBufferPool = er.New("closed buffer pool")

// ErrClosedBuffer 是表示缓冲器已关闭的错误的变量。
var ErrClosedBuffer = er.New("closed buffer")

// FIFO缓冲器的接口
type Buffer interface {
	// 获取本缓冲器的容量
	Cap() uint32
	// 获取缓冲器的数据数量
	Len() uint32
	// 向缓冲器中放入数据，若缓冲器已关闭则直接返回非nil错误值
	Put(datum interface{}) (bool, error)
	// 从缓冲器中获取数据，若已关闭则返回错误值
	Get() (interface{}, error)
	// 关闭缓冲器
	Close() bool
	// 判断缓冲器是否已关闭
	Closed() bool
}

type myBuffer struct {
	// 存放数据的通道
	ch chan interface{}
	// 缓冲器的关闭状态，0：未关闭，1：已关闭
	closed uint32
	// 关闭缓冲器时的竞态锁
	closingLock sync.RWMutex
}

func (buf *myBuffer) Cap() uint32 {
	return uint32(cap(buf.ch))
}

func (buf *myBuffer) Len() uint32 {
	return uint32(len(buf.ch))
}

func (buf *myBuffer) Put(datum interface{}) (bool, error) {
	buf.closingLock.RLock()
	defer buf.closingLock.RUnlock()

	if buf.Closed() {
		return false, ErrClosedBuffer
	}
	var ok bool
	select {
	case buf.ch <- datum:
		ok = true
	default:
		ok = false
	}
	return ok, nil
}

func (buf *myBuffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-buf.ch:
		if !ok {
			return nil, ErrClosedBuffer
		}
		return datum, nil
	default:
		return nil, nil
	}
}

func (buf *myBuffer) Close() bool {
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) {
		buf.closingLock.Lock()
		defer buf.closingLock.Unlock()
		close(buf.ch)
		return true
	}
	return false
}

func (buf *myBuffer) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}

// 创建一个缓冲器，size为缓冲器的大小
func NewBuffer(size uint32) (Buffer, error) {
	if size == 0 {
		errMsg := fmt.Sprintf("illegal size for buffer: %d", size)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	return &myBuffer{
		ch: make(chan interface{}, size),
	}, nil
}
