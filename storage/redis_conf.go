package storage

import (
	"github.com/go-redis/redis"
)

const (
	redisAddr = "172.17.0.1:6379"
	redisPwd  = ""
)

//NewRedisClient init redis conf
func NewRedisClient() *redis.Client {
	ops := new(redis.Options)
	ops.Addr = redisAddr
	ops.DB = 0
	ops.Password = redisPwd
	return redis.NewClient(ops)
}
