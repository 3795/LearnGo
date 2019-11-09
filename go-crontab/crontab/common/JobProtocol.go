package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// 定时任务
type Job struct {
	Name     string `json:"name"`     // 任务名称
	Command  string `json:"command"`  // shell命令
	CronExpr string `json:"cronExpr"` // cron表达式
}

// 任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 // 调度任务
	Expr     *cronexpr.Expression // 解析好的cron表达式
	NextTime time.Time            // 下次调度时间
}

// 任务执行状态
type JobExecuteInfo struct {
	Job        *Job               // 任务信息
	PlanTime   time.Time          // 理论上的调度时间
	RealTime   time.Time          // 实际上的调度时间
	CancelCtx  context.Context    // 任务的context
	CancelFunc context.CancelFunc // 取消命令执行的函数
}

type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo // 执行状态
	Output      []byte          // 脚本输出结果
	Err         error           // 脚本错误原因
	StartTime   time.Time       // 启动时间
	EndTime     time.Time       // 结束时间
}

// 变化事件
type JobEvent struct {
	EventType int // 变化类型，如新增，删除等
	Job       *Job
}

// 反序列化Job结构
func UnpackJob(value []byte) (ret *Job, err error) {
	job := &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		ret = nil
		return
	}

	return
}

// 构建任务变化事件
func BuildJobEvent(jobEventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: jobEventType,
		Job:       job,
	}
}

// 构建任务执行计划
func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)

	// 解析Job的cron表达式
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return
}

// 构造执行状态信息
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime,
		RealTime: time.Now(),
	}

	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}

// 提取worker节点的IP
func ExtracWorkerIP(regKey string) string {
	return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}
