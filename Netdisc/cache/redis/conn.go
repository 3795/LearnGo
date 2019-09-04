package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"sync"
	"time"
)

var (
	pool     *redis.Pool
	host     = "127.0.0.1"
	port     = 6379
	password = "123456"
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			// 建立连接
			address := fmt.Sprintf("%s:%d", host, port)
			conn, e = redis.Dial("tcp", address)
			if e != nil {
				log.Println(e)
				return nil, e
			}

			// 使用密码进行访问
			if _, e = conn.Do("AUTH", password); e != nil {
				log.Println(e)
				_ = conn.Close()
				return nil, e
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
	}
}

func RedisPool() *redis.Pool {
	if pool == nil {
		o := sync.Once{}
		o.Do(func() {
			pool = newRedisPool()
		})
	}
	return pool
}
