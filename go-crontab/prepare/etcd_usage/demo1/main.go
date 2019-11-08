package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		//lease   clientv3.Lease
		//watcher clientv3.Watcher
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"192.168.56.101:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)

	// 放入KV
	Put(kv)
	// 获取KV
	//Get(kv)
	// 删除KV
	//Del(kv)

	// 申请一个lease
	//lease = clientv3.NewLease(client)
	//Lease(lease, kv)

	//watcher = clientv3.NewWatcher(client)
	//Watch(kv, watcher)

	//Operation(kv)
}

/**
OP操作
*/
func Operation(kv clientv3.KV) {
	var (
		putOp  clientv3.Op
		getOp  clientv3.Op
		opResp clientv3.OpResponse
		err    error
	)

	// 执行Put操作
	putOp = clientv3.OpPut("/cron/jobs/job8", "12334566")
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入Revision：", opResp.Put().Header.Revision)

	// 执行Get操作
	getOp = clientv3.OpGet("/cron/jobs/job8")
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("数据Revision：", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据Value：", string(opResp.Get().Kvs[0].Value))
}

/**
监听功能
*/
func Watch(kv clientv3.KV, watcher clientv3.Watcher) {
	var (
		getResp            *clientv3.GetResponse
		err                error
		watchStartRevision int64
		watchRespChan      <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
	)
	// 持续对值做出改变
	go func() {
		for {
			_, _ = kv.Put(context.TODO(), "/cron/jobs/job7", "i am job7")
			_, _ = kv.Delete(context.TODO(), "/cron/jobs/job7")
			time.Sleep(1 * time.Second)

		}
	}()

	// 先Get到当前的值，并监听后续变化
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}

	// key如果存在
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值：", string(getResp.Kvs[0].Value))
	}

	// 当前etcd集群事务ID，是单调递增的
	watchStartRevision = getResp.Header.Revision + 1

	// 启动监听
	fmt.Println("从该版本开始监听: ", watchStartRevision)

	ctx, cancelFunc := context.WithCancel(context.Background())
	// 5秒后终止监听
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})

	// 监听值的变化
	watchRespChan = watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "Revision：", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了：", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}

/**
租约功能
*/
func Lease(lease clientv3.Lease, kv clientv3.KV) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		err            error
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
	)
	// 申请一个10秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	// 获取租约ID
	leaseId = leaseGrantResp.ID

	// 这样设置后，每过3秒左右，就会续租租约，重置过期时间
	//if keepRespChan, err = lease.KeepAlive(context.Background(), leaseId); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 使用timeOut型context，在5秒的时间内，还是会多次续租租约，重置过期时间
	// 超过5秒后，不再续租租约，KV到租约有效期后，就自动过期了
	timeOutContext, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if keepRespChan, err = lease.KeepAlive(timeOutContext, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	// 处理续约应答协程
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					fmt.Println("已取消自动续约，续约即将过期")
					goto END
				} else {
					fmt.Println("收到自动续约应答：", keepResp.ID)
				}
			}
		}
	END:
	}()

	// Put一个KV，并与租约相关联，实现10秒后过期
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功：", putResp.Header.Revision)

	// 定时查看KV是否过期
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("租约已过期")
			break
		}
		fmt.Println("租约暂未过期", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
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
