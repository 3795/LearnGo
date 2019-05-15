package datasource

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"project/lottery/project/conf"
	"sync"
	"time"
)

// 连接到Redis
var redisLock sync.Mutex
var cacheInstance *RedisConn

type RedisConn struct {
	pool 	*redis.Pool
	showDebug 	bool
}

func (r *RedisConn) Do(command string, args ...interface{}) (response interface{}, err error) {
	conn := r.pool.Get()
	defer conn.Close()

	t1 := time.Now().UnixNano()
	response, err = conn.Do(command, args...)
	if err != nil {
		e := conn.Err()
		if e != nil {
			log.Println("redishelper.Do error:", err, e)
		}
	}
	t2 := time.Now().UnixNano()
	if r.showDebug {
		fmt.Printf("[redis] [info] [%dus]cmd=%s, err=%s, args=%v, reply=%s\n", (t2-t1)/1000, command, err, args, response)
	}
	return response, err
}

func (r *RedisConn) ShowDebug(b bool) {
	r.showDebug = b
}

func InstanceCache() *RedisConn {
	if cacheInstance != nil {
		return cacheInstance
	}
	redisLock.Lock()
	defer redisLock.Unlock()

	if cacheInstance != nil {
		return cacheInstance
	}
	return NewCache()
}

func NewCache() *RedisConn {
	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.RdsCache.Host, conf.RdsCache.Port))
			if err != nil {
				log.Fatal("redishelper.NewCache Dial error: ", err)
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle: 	10000,
		MaxActive: 	10000,
		IdleTimeout:	0,
		Wait:	false,
		MaxConnLifetime:	0,
	}

	instance := &RedisConn{
		pool:	&pool,
	}
	cacheInstance = instance
	cacheInstance.ShowDebug(true)
	return cacheInstance
}