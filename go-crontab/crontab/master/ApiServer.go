package master

import (
	"net"
	"net/http"
)

// 任务的HTTP接口
type ApiServer struct {
	httpServer *http.Server
}

// 创建单例对象
var (
	G_apiServer *http.Server
)

func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpServer    *http.Server
		staticDir     http.Dir     // 静态文件根目录
		staticHandler http.Handler // 静态文件的HTTP回调
	)
	// 配置路由

	// 静态文件目录
	//staticDir = http.Dir(G_C)

	// 启动TCP监听
	//if listener, err = net.Listen("tcp")
	return
}
