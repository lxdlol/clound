package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var(
	pool   *redis.Pool
	redispass  ="123456"
	redishost ="127.0.0.1:6739"
)



func newRedisPool()*redis.Pool{
	return &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			conn, e = redis.Dial("tcp", redishost)
			if e!=nil{
				fmt.Println(e)
				return nil,e
			}
			if _, e := conn.Do("AUTH", redispass);e!=nil{
				fmt.Println(e)
				return nil,e
			}
			return conn,nil
		},
		MaxIdle:50,
		MaxActive:30,
		IdleTimeout:300 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func init(){
	pool=newRedisPool()
}

func RedisPool()*redis.Pool{
	return pool
}