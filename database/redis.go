package data

import (
	"context"
	"github.com/redis/go-redis/v9"
	c "pigpq/config"
	"pigpq/global"
)

func initRedis() {
	global.RedisClient = redis.NewClient(&redis.Options{
		Addr:     c.Config.Redis.Host + ":" + c.Config.Redis.Port,
		Password: c.Config.Redis.Password,
		DB:       c.Config.Redis.Database,
	})
	var ctx = context.Background()

	_, err := global.RedisClient.Ping(ctx).Result()

	if err != nil {
		panic("Redis connection failed：" + err.Error())
	}
}
