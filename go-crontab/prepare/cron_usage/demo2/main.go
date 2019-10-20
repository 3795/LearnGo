package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {
	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob
	)

	scheduleTable = make(map[string]*CronJob)
	now = time.Now()

	// 第一个定时任务
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	scheduleTable["job1"] = cronJob

	// 第二个定时任务
	expr = cronexpr.MustParse("*/1 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	scheduleTable["job2"] = cronJob

	// 启动监听协程
	go func() {
		for {
			now := time.Now()
			for jobName, cronJob := range scheduleTable {
				// 判断是否过期
				if cronJob.nextTime.Equal(now) || cronJob.nextTime.Before(now) {
					go func(jobName string) {
						fmt.Println("执行了 ", jobName)
					}(jobName)
					cronJob.nextTime = cronJob.expr.Next(now)
				}
			}

			// 睡眠100ms，防止占用CPU
			time.Sleep(100)
		}
	}()

	time.Sleep(100 * time.Second)
}
