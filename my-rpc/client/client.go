package client

import (
	"context"
	"errors"
)

var ErrorShutdown = errors.New("Client is shut down")

type Call struct {
	ServiceMethod string	// 服务名，方法名
	Args	interface{}		// 参数
	Reply	interface{}		// 返回值（指针类型）
	Error	error		// 错误信息
	Done	chan *Call		// 在调用结束时激活
}

type simpleClient struct {
	//codec
}

type RPCClient interface {
	Go(ctx context.Context, serviceMethod string, arg interface{}, reply interface{}, done chan *Call) *Call
	Call(ctx context.Context, serviceMethod string, arg interface{}, reply interface{}) error
	Close() error
	IsShutDown() bool
}

