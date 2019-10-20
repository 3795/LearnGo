package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// 日志记录模型
type LogRecord struct {
	JobName   string    `bson:"jobName"`   // 任务名
	Command   string    `bson:"command"`   // shell命令内容
	Err       string    `bson:"err"`       // 脚本错误信息
	Content   string    `bson:"content"`   // 脚本输出内容
	TimePoint TimePoint `bson:"timePoint"` // 执行时间点
}

// 过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

// startTime小于某时间
// {"$lt": timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// {"timePoint.startTime": {"$lt": timestamp}}
type DeleteCond struct {
	before TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
	)

	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://127.0.0.1:27017")); err != nil {
		fmt.Println(err)
		return
	}

	database = client.Database("cron")

	collection = database.Collection("log")

	//Find(collection)
	Delete(collection)

}

func Delete(collection *mongo.Collection) {
	var (
		delResult *mongo.DeleteResult
		err       error
	)
	// 构建删除条件
	delCond := &DeleteCond{
		before: TimeBeforeCond{
			Before: time.Now().Unix(),
		},
	}

	// 执行删除

	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除的行数: ", delResult.DeletedCount)
}

/**
查询记录
*/
func Find(collection *mongo.Collection) {
	var (
		cond   *FindByJobName
		cursor *mongo.Cursor
		record *LogRecord
		err    error
	)
	cond = &FindByJobName{JobName: "jobName"}

	if cursor, err = collection.Find(context.TODO(), cond, options.Find().SetSkip(0), options.Find().SetLimit(2)); err != nil {
		fmt.Println(err)
		return
	}

	defer cursor.Close(context.TODO())

	// 遍历查询结果
	for cursor.Next(context.TODO()) {
		record = &LogRecord{}

		// 反序列化
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(*record)
	}

}

/**
插入一条记录
*/
func InsertOne(collection *mongo.Collection) {
	record := &LogRecord{
		JobName: "jobName",
		Command: "echo hello",
		Err:     "",
		Content: "",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	var (
		result *mongo.InsertOneResult
		err    error
	)

	// 插入记录
	if result, err = collection.InsertOne(context.Background(), record); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("生成的ID", result.InsertedID)
}
