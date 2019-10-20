package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

var (
	expr     *cronexpr.Expression
	err      error
	now      time.Time
	nextTime time.Time
)

func main() {
	// 新建cron表达式
	if expr, err = cronexpr.Parse("*/1 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	// 获取当前时间
	now = time.Now()
	// 计算下一个调度时间点
	nextTime = expr.Next(now)
	ch := make(chan int, 1)
	go func() {
		for {
			time.AfterFunc(nextTime.Sub(now), func() {
				fmt.Println("被调度了一次：", time.Now().String())
				ch <- 1
			})
			now = nextTime
			nextTime = expr.Next(now)
			<-ch
		}
	}()

	//// 获取当前时间
	//now = time.Now()
	//// 计算下一个调度时间点
	//nextTime = expr.Next(now)
	//
	//time.AfterFunc(nextTime.Sub(now), func() {
	//	fmt.Println("被调度了：", nextTime)
	//})

	time.Sleep(600 * time.Second)
	fmt.Println("执行完成")

}
