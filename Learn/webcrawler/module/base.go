package module

import "net/http"

// 汇集组件内部计数的类型
type Counts struct {
	// 调用的次数
	CalledCount uint64

	// 接受调用次数
	AcceptedCount uint64
	// 调用成功完成次数
	CompletedCount uint64
	// 实时处理数
	HandlingNumber uint64
}

// 组件摘要结构的类型
type SummaryStrut struct {
	ID        MID         `json:"id"`
	Called    uint64      `json:"called"`
	Accepted  uint64      `json:"accepted"`
	Completed uint64      `json:"completed"`
	Handling  uint64      `json:"handling"`
	Extra     interface{} `json:"extra, omitempty"`
}

// 组件的基础接口类型，该接口的实现类型必须是并发安全的
type Module interface {
	ID() MID                         // 当前组件的ID
	Addr() string                    // 当前组件的网络地址
	Score() uint64                   // 当前组件的评分
	SetScore(score uint64)           // 设置当前组件的评分
	ScoreCalculator() CalculateScore // 获取评分计算器
	CalledCount() uint64             // 获取当前组件被调用的次数
	AcceptedCount() uint64           // 获取当前组件接受的调用次数，当组件超负载时会拒绝调用
	CompletedCount() uint64          // 组件已成功完成的调用计数
	HandlingNumber() uint64          // 当前组件正在处理的调用的数量
	Counts() Counts                  // 用于一次性获取所有计数
	Summary() SummaryStrut           // 获取组件摘要
}

// Downloader 代表下载器的接口类型。
// 该接口的实现类型必须是并发安全的！
type Downloader interface {
	Module
	// 根据请求获取内容并返回响应
	Download(req *Request) (*Response, error)
}

// Analyzer 代表分析器的接口类型。
// 该接口的实现类型必须是并发安全的！
type Analyzer interface {
	Module
	// 返回当前分析器使用的响应解析函数的列表
	RespParsers() []ParseResponse
	// Analyze根据规则分析响应并返回请求的和条目
	// 响应需要经过多个响应解析函数的处理，然后合并结果
	Analyze(resp *Response) ([]Data, []error)
}

// 用于解析Http响应的函数的类型
type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]Data, []error)

type Pipeline interface {
	Module
	// 返回当前条目处理管道使用的条目处理函数列表
	ItemProcessors() []ProcessItem
	// 向条目处理管道发送条目，条目会经过多个处理函数的处理
	Send(item Item) []error
	// 该方法返回一个布尔值，该值表示当前条目处理管道是否是快速失败的
	// 这里的快速失败是指：只要在处理条目的过程中某一个步骤上出错，那么条目处理管道就会停止执行并报告错误
	FailFast() bool
	// 设置其是否是快速失败的
	SetFailFast(failFast bool)
}

// 用于处理条目的函数的类型
type ProcessItem func(item Item) (result Item, err error)
