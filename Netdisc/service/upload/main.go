package main

import (
	"LearnGo/Netdisc/handler"
	"fmt"
	"net/http"
)

const (
	address = "0.0.0.0:8080"
)

func main() {
	// 处理静态资源
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("Netdisc/static"))))

	// 上传相关接口
	http.HandleFunc("/file/upload", handler.UploadHandler)

	// 用户相关接口
	http.HandleFunc("/", handler.SignInHandler)
	http.HandleFunc("/user/signup", handler.SignInHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	fmt.Printf("上传服务启动中，开始监听 %s \n", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}

}
