package ioc

import (
	"github.com/redis/go-redis/v9"

	"geektime/webook/config"
)

func InitRedis() redis.Cmdable {
	rCfg := config.Config.Redis
	return redis.NewClient(&redis.Options{
		Addr: rCfg.Addr,
	})
}
