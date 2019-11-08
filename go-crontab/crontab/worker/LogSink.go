package worker

import (
	"LearnGo/go-crontab/crontab/common"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	G_logSink *LogSink
)

// 批量写入日志
func (logSink *LogSink) saveLogs(batch *common.LogBatch) {
	_, err := logSink.logCollection.InsertMany(context.TODO(), batch.Logs)
	if err != nil {
		fmt.Printf("写入日志失败, %v\n", err.Error())
	}
}

// 日志存储携程，将日志加入日志批次，批次到达一定时间后一次性向MongoDB写入
func (logSink *LogSink) writeLoop() {
	var (
		log          *common.JobLog
		logBatch     *common.LogBatch // 当前的批次
		commitTimer  *time.Timer      // 提交的时间
		timeoutBatch *common.LogBatch // 超时批次
	)

	for {
		select {
		case log = <-logSink.logChan:
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				// 设置该批次自动提交的时间
				commitTimer = time.AfterFunc(time.Duration(G_config.JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() {
							logSink.autoCommitChan <- batch
						}
					}(logBatch))
			}

			// 将日志追加到批次中
			logBatch.Logs = append(logBatch.Logs, log)

			// 如果批次满了，就立即发送
			if len(logBatch.Logs) >= G_config.JobLogBatchSize {
				// 发送日志
				logSink.saveLogs(logBatch)
				logBatch = nil
				// 取消定时器
				commitTimer.Stop()
			}
		// 超时自动提交机制
		case timeoutBatch = <-logSink.autoCommitChan:
			if timeoutBatch != logBatch {
				/**
				场景：
				1. 日志批次到达大小上限，进入提交逻辑
				2. 计时器刚好触发自动提交逻辑，向通道中写入logBatch
				3. 大小上限逻辑写入日志成功，将logBatch置为nil
				4. 如果没有这个判断机制，那么自动提交逻辑可能会重复提交日志
				*/
				continue // 跳过已提交的日志批次
			}
			// 将该批次日志写入
			logSink.saveLogs(timeoutBatch)
			logBatch = nil

		}
	}
}

func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)

	// 建立MongoDB连接
	if client, err = mongo.Connect(context.TODO(),
		options.Client().ApplyURI(G_config.MongodbUri),
		options.Client().SetConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond)); err != nil {
		fmt.Printf("连接MongoDB失败，%v\n", err.Error())
		return
	}

	// 选择db和collection
	G_logSink = &LogSink{
		client:         client,
		logCollection:  client.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	// 启动日志处理携程
	go G_logSink.writeLoop()
	return
}

// 发送日志
func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	// 先尝试将该日志加入通道
	case logSink.logChan <- jobLog:
	// 如果通道满了，就直接丢弃该日志，避免程序被通道阻塞住
	default:
	}
}
