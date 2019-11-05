package worker

import "LearnGo/go-crontab/crontab/common"

// 任务调度
type Scheduler struct {
	jobEventChan      chan *common.JobEvent              // 任务事件队列
	jobPlanTable      map[string]*common.JobSchedulePlan // 任务调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo  // 任务执行表
	jobResultChan     chan *common.JobExecuteResult      // 任务结果队列
}

var (
	G_scheduler *Scheduler
)
