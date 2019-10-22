package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"192.168.74.128:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)

	Put(kv)
	//Get(kv)
	Del(kv)
}

func Del(kv clientv3.KV) {
	var (
		delResp *clientv3.DeleteResponse
		err     error
	)

	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job2"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("本次删除的个数为%d\n", delResp.Deleted)
		fmt.Println(len(delResp.PrevKvs))
		if len(delResp.PrevKvs) > 0 {
			for _, kvPair := range delResp.PrevKvs {
				fmt.Printf("key = %v, value = %v\n", kvPair.Key, kvPair.Value)
			}
		}
	}
}

func Get(kv clientv3.KV) {
	var (
		getResp *clientv3.GetResponse
		err     error
	)

	if getResp, err = kv.Get(context.TODO(), "/cron/jobs", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(getResp.Kvs)
	}
}

/**
放入KV
*/
func Put(kv clientv3.KV) {

	var (
		putResp *clientv3.PutResponse
		err     error
	)

	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job2", "world"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(putResp.Header.Revision)
	}
}
