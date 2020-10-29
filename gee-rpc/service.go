package gee_rpc

import (
	"LearnGo/gee-rpc/codec"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int        // 标记这是一个gee-rpc请求
	CodecType   codec.Type // 选择不同的编码方式
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

type Server struct{}

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
	// TODO 此处暂时还不知道请求参数的类型
	// 暂时只支持string类型的参数
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
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
	// TODO should call registered rpc methods to get the right replyv
	// 暂时只打印请求参数并返回一个hello world信息
	defer wg.Done()
	log.Print("Handle Request: ", req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("gee-rpc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}
