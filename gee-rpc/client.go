package gee_rpc

import (
	"LearnGo/gee-rpc/codec"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{} // 方法调用的参数
	Reply         interface{} // 方法调用的返回类型
	Error         error       // 远程调用返回的错误信息
	Done          chan *Call  // 回调通道
}

func (call *Call) done() {
	call.Done <- call
}

var ErrShutdown = errors.New("connection is shut down")

type Client struct {
	cc       codec.Codec // 消息的编解码器
	opt      *Option     // 告诉Server端消息的编码格式
	sending  sync.Mutex  // 互斥锁，保证消息的有序发送，防止多个请求报文混淆
	header   codec.Header
	mu       sync.Mutex
	seq      uint64           // 发送的请求编号，每个请求有唯一的编号
	pending  map[uint64]*Call // 存储为处理完的请求，键是编号，值是Call实例
	closing  bool             // 用户主动关闭时，调用Close方法
	shutdown bool             // 当有错误发生时，shutdown值置为true
}

func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.cc.Close()
}

func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.shutdown && !client.closing
}

// 将call添加到client.pending中，并更新client.seq
func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}

// 从pending中删除call
func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

// 服务端或客户端发生错误时，将shutdown置为true，并将错误信息通知所有pending状态的call
func (client *Client) terminateCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()
	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

// 客户端接收响应
func (client *Client) receive() {
	var err error
	for err == nil {
		// 读取header时就出错，说明这个请求有问题
		var h codec.Header
		if err = client.cc.ReadHeader(&h); err != nil {
			break
		}
		call := client.removeCall(h.Seq) // 从pending状态删除
		switch {
		case call == nil: // 正常情况下call不会为空
			err = client.cc.ReadBody(nil)
		case h.Error != "": // 服务端出错了
			call.Error = fmt.Errorf(h.Error)
			err = client.cc.ReadBody(nil)
			call.done()
		default: // 正常处理
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	// 当错误发生时，关闭所有在等待队列中的call，上面是一个死循环，不发生错误到不了这一步
	client.terminateCalls(err)
}

func NewClient(conn net.Conn, opt *Option) (*Client, error) {
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("rpc client: codec error:", err)
		return nil, err
	}

	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error:", err)
		_ = conn.Close()
		return nil, err
	}
	return newClientCodec(f(conn), opt), nil
}

func newClientCodec(cc codec.Codec, opt *Option) *Client {
	client := &Client{
		opt:     opt,
		cc:      cc,
		seq:     1, // 序列号从1开始
		pending: make(map[uint64]*Call),
	}
	go client.receive()
	return client
}

func parseOptions(opts ...*Option) (*Option, error) {
	if len(opts) == 0 || opts[0] == nil {
		return DefaultOption, nil
	}

	if len(opts) == 1 {
		return nil, errors.New("number of options is more than 1")
	}
	opt := opts[0]
	opt.MagicNumber = DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = DefaultOption.CodecType
	}
	return opt, nil
}

// 向Server端发送一个连接
func Dial(network string, address string, opts ...*Option) (*Client, error) {
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	client, err := NewClient(conn, opt)
	defer func(client *Client) {
		if client == nil {
			_ = conn.Close()
		}
	}(client)

	return client, err
}

// Client向Server发送消息
func (client *Client) send(call *Call) {
	client.sending.Lock()
	defer client.sending.Unlock()

	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	client.header.ServiceMethod = call.ServiceMethod
	client.header.Seq = seq
	client.header.Error = ""

	// 如果调用出现了问题
	if err := client.cc.Write(&client.header, call.Args); err != nil {
		call := client.removeCall(seq) // 从待执行队列中删除这个请求
		if call != nil {
			call.Error = err // 填充请求调用的错误信息
			call.done()      // 将请求置为已完成状态
		}
	}
	// 如果调用没有问题，就不调用call.done()，等待服务端响应，在receive函数中将这个调用置为完成状态
}

// 异步进行远程调用，返回一个Call实例
func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered") // 这个通道的容量必须大于0
	}

	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	client.send(call)
	return call
}

// 远程调用的回调函数，是对Go函数的封装，阻塞call.Done，等待响应返回，是一个同步接口
func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}
