package data

import (
	c "pigpq/config"
	"sync"
)

var once sync.Once

func InitDatabase() {
	once.Do(func() {
		if c.Config.Mysql.Enable {
			// 初始化 mysql
			initMysql()
		}

		if c.Config.Redis.Enable {
			// 初始化 redis
			initRedis()
		}
	})
}
