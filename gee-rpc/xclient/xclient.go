package xclient

import (
	gee_rpc "LearnGo/gee-rpc"
	"context"
	"reflect"
	"sync"
)

type XClient struct {
	d       Discovery
	mode    SelectMode
	opt     *gee_rpc.Option
	mu      sync.Mutex
	clients map[string]*gee_rpc.Client
}

func NewXClient(d Discovery, mode SelectMode, opt *gee_rpc.Option) *XClient {
	return &XClient{
		d:       d,
		mode:    mode,
		opt:     opt,
		clients: make(map[string]*gee_rpc.Client),
	}
}

// 关闭已建立的连接
func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	for key, client := range xc.clients {
		_ = client.Close()
		delete(xc.clients, key)
	}
	return nil
}

func (xc *XClient) dial(rpcAddr string) (*gee_rpc.Client, error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	client, ok := xc.clients[rpcAddr] // 有缓存就先取缓存中的客户端
	if ok && !client.IsAvailable() {
		// 如果从数组中取出来的客户端不可用
		_ = client.Close()
		delete(xc.clients, rpcAddr)
		client = nil
	}

	if client == nil {
		var err error
		client, err = gee_rpc.XDial(rpcAddr, xc.opt) // 建立连接
		if err != nil {
			return nil, err
		}
		xc.clients[rpcAddr] = client // 将建立的连接保存起来准备复用
	}
	return client, nil
}

func (xc *XClient) call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
	client, err := xc.dial(rpcAddr)
	if err != nil {
		return err
	}
	return client.Call(ctx, serviceMethod, args, reply)
}

func (xc *XClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}
	return xc.call(rpcAddr, ctx, serviceMethod, args, reply)
}

func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	servers, err := xc.d.GetAll()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var e error
	replyDone := reply == nil
	ctx, cancel := context.WithCancel(ctx)
	for _, rpcAddr := range servers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var clonedReply interface{}
			if reply != nil {
				clonedReply = reflect.New(reflect.ValueOf(reply).Elem().Type())
			}
			err := xc.call(rpcAddr, ctx, serviceMethod, args, clonedReply) // 广播注册方法
			mu.Lock()
			if err != nil && e == nil {
				e = err
				cancel() // 如果调用失败，则context直接失败
			}
			if err == nil && !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(clonedReply).Elem())
				replyDone = true
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	return e
}
