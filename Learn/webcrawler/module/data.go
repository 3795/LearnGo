package module

import "net/http"

// 数据的接口类型
type Data interface {
	// 判断数据是否有效
	Valid() bool
}

// 数据请求的类型
type Request struct {
	// HTTP请求
	httpReq *http.Request
	// 请求的深度
	depth uint32
}

// 新建一个请求实例
func NewRequest(httpReq *http.Request, depth uint32) *Request {
	return &Request{httpReq, depth}
}

// 获取HTTP请求
func (req *Request) HttpReq() *http.Request {
	return req.httpReq
}

// 获取请求的深度
func (req *Request) Depth() uint32 {
	return req.depth
}

// 验证该请求是否有效
func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

type Response struct {
	// HTTP响应
	httpResp *http.Response
	// 响应的深度
	depth uint32
}

// 新建一个响应实例
func NewResponse(httpResp *http.Response, depth uint32) *Response {
	return &Response{httpResp, depth}
}

func (resp *Response) HttpResp() *http.Response {
	return resp.httpResp
}

func (resp *Response) Depth() uint32 {
	return resp.depth
}

// 验证响应是否有效
func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

// 条目的类型
type Item map[string]interface{}

// 验证条目是否有效
func (item Item) Valid() bool {
	return item != nil
}
