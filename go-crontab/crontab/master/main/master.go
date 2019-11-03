package main

import (
	"LearnGo/go-crontab/crontab/master"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
	flag.Parse()
}

// 初始化线程数量,指定线程数量为CPU核数
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	initArgs()

	initEnv()

	// 加载配置
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 启动Api HTTP服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	// 启动Etcd服务
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(5 * time.Second)
	}

ERR:
	fmt.Println(err)
}
