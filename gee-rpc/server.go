package gee_rpc

import (
	"LearnGo/gee-rpc/codec"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

const MagicNumber = 0x3bef5c

const (
	connected        = "200 Connected to GeeRPC"
	defaultRPCPath   = "/_geerpc_"
	defaultDebugPath = "/debug/geerpc"
)

type Option struct {
	MagicNumber       int           // 标记这是一个gee-rpc请求
	CodecType         codec.Type    // 选择不同的编码方式
	ConnectionTimeout time.Duration // 连接超时时间，超时时间设置为0时，表示不对超时时间进行限制
	HandleTimeout     time.Duration // 处理Request时的超时时间
}

var DefaultOption = &Option{
	MagicNumber:       MagicNumber,
	CodecType:         codec.GobType,
	ConnectionTimeout: 10 * time.Second,
}

type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

// 接收链接
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error: ", err)
			return
		}

		go server.ServerConn(conn)
	}
}

// 处理链接
func (server *Server) ServerConn(conn io.ReadWriteCloser) {
	defer func() {
		_ = conn.Close()
	}()

	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}

	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invaild magic number %x", opt.MagicNumber)
		return
	}

	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	server.serveCodec(f(conn))
}

var invalidRequest = struct{}{}

// 处理编码信息
func (server *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex) // 加锁，确保每次发送的是一个完成的Response
	wg := new(sync.WaitGroup)  // wait until all request are handled

	// 一个网络连接允许接收多个请求，所以用for无限制的等待请求的到来
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break // 当解析header失败时，就直接关闭连接
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending) // 直接返回错误信息
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

// request stores all information of a call（该结构保存一个远程调用的所有信息，包括函数的参数和返回值）
type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
	mtype        *methodType
	svc          *service
}

// 读取请求的头信息
func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	h := &codec.Header{}
	if err := cc.ReadHeader(h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error: ", err)
		}
		return nil, err
	}
	return h, nil
}

// 读取请求本身的内容
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	req.svc, req.mtype, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}
	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReply()

	// 确保argvi是一个指针类型参数，ReadBody需要一个指针类型的参数
	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}

	if err = cc.ReadBody(argvi); err != nil {
		log.Println("rpc server: read argv err: ", err)
	}
	return req, nil
}

// 向客户端发送响应
func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	// 可以并发处理请求，但是给客户端返回response必须要一个一个返回，以免多个response交织在一起
	// 所以这个地方用一个锁来处理
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error: ", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	timeout := 5 * time.Second // 设置处理请求超时时间为5秒
	defer wg.Done()

	called := make(chan struct{})
	sent := make(chan struct{})

	go func() {
		err := req.svc.Call(req.mtype, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
		}
		server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		sent <- struct{}{}
	}()

	if timeout == 0 {
		<-called
		<-sent
		return
	}

	select {
	case <-time.After(timeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", timeout)
		server.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called:
		<-sent
	}
}

func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}

// 将服务注册到server
func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

// 根据服务名找到service
func (server *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	svc = svci.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}
	return
}
