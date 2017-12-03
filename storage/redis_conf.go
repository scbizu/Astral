package storage

import (
	"github.com/go-redis/redis"
)

const (
	redisAddr = "localhost:32826"
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
