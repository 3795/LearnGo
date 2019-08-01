package scheduler

// 代表参数容器的接口类型
type Args interface {
	// 自检参数有效性
	Check() error
}

// 请求相关的参数容器类型
type RequestArgs struct {
	// 可以接受的ULR的主域名的列表，URL主域名不在列表中的会被忽略
	AcceptedDomains []string `json:"accepted_domains"`
	// 爬取的最大深度，超过此深度则放弃继续爬取
	MaxDepth uint32 `json:"max_depth"`
}

func (args *RequestArgs) Check() error {
	if args.AcceptedDomains == nil {
		return genError("nil accepted primary domain list")
	}
	return nil
}
