package main

import (
	gee_rpc "LearnGo/gee-rpc"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	addr := make(chan string)
	go startServer(addr) // 启动服务器

	client, _ := gee_rpc.Dial("tcp", <-addr)
	defer func() {
		_ = client.Close()
	}()

	time.Sleep(time.Second) //睡1秒，给服务器时间

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("gee-rpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error", err)
			}
			log.Println("reply: ", reply)
		}(i)
	}
	wg.Wait()
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
