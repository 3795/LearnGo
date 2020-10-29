package main

import (
	gee_rpc "LearnGo/gee-rpc"
	"LearnGo/gee-rpc/codec"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	addr := make(chan string)
	go startServer(addr) // 启动服务器

	conn, _ := net.Dial("tcp", <-addr) // 建立自己到自己的连接
	defer func() {
		_ = conn.Close()
	}()

	time.Sleep(time.Second) //睡1秒，给服务器时间

	// 发送测试请求
	_ = json.NewEncoder(conn).Encode(gee_rpc.DefaultOption)
	cc := codec.NewGobCodec(conn)

	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}

		// 向服务端发送信息
		_ = cc.Write(h, fmt.Sprintf("gee-rpc %d", h.Seq))

		// 接收从服务端返回的信息
		var reply string
		_ = cc.ReadHeader(h)
		_ = cc.ReadBody(&reply)
		log.Println("receive reply:", reply)
	}

}

func startServer(addr chan string) {
	l, err := net.Listen("tcp", "127.0.0.1:7000")
	if err != nil {
		log.Fatal("network error: ", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	gee_rpc.Accept(l)
}
