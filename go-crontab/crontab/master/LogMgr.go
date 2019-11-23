package master

import (
	"LearnGo/go-crontab/crontab/common"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 日志管理
type LogMgr struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

func InitLogMgr() (err error) {
	var (
		client *mongo.Client
	)

	// 建立MongoDB的连接
	if client, err = mongo.Connect(context.TODO(),
		options.Client().ApplyURI(G_config.MongodbUri),
		options.Client().SetConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond)); err != nil {
		fmt.Printf("连接MongoDB失败，%v\n", err.Error())
		return
	}

	G_logMgr = &LogMgr{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}

	return
}

// 查看任务日志
func (logMgr *LogMgr) ListLog(name string, skip int, limit int) (logArr []*common.JobLog, err error) {
	var (
		filter  *common.JobLogFilter
		logSort *common.SortLogByStartTime
		cursor  *mongo.Cursor
		jobLog  *common.JobLog
	)

	logArr = make([]*common.JobLog, 0)
	// 设置过滤条件
	filter = &common.JobLogFilter{JobName: name}
	// 按照时间顺序倒序排列
	logSort = &common.SortLogByStartTime{SortOrder: -1}
	// 查询
	s := int64(skip)
	l := int64(limit)
	if cursor, err = logMgr.logCollection.Find(context.TODO(), filter,
		&options.FindOptions{Sort: logSort, Skip: &s, Limit: &l}); err != nil {
		return
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}

		// 反序列化BSON
		if err = cursor.Decode(jobLog); err != nil {
			// 忽略不合格式的日志
			continue
		}
		logArr = append(logArr, jobLog)
	}
	return
}
