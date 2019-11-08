package worker

import (
	"LearnGo/go-crontab/crontab/common"
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// 任务管理器
type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	// 单例模式
	G_jobMgr *JobMgr
)

// 监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)

	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	// 遍历当前事件
	for _, kvpair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			// 同步给调度携程
			G_scheduler.PushJobEvent(jobEvent)
		}
	}

	// 从该revision向后监听变化事件
	go func() {
		watchStartRevision = getResp.Header.Revision + 1
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:

				}
			}
		}
	}()

	return
}

// 创建任务执行锁
func (jobMgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return
}
