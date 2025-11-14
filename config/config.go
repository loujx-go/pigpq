package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	. "pigpq/config/autoload"
	"pigpq/pkg/untils"
)

// Conf 配置项主结构体
type Conf struct {
	AppConfig `mapstructure:"app"`
	Mysql     MysqlConfig  `mapstructure:"mysql"`
	Redis     RedisConfig  `mapstructure:"redis"`
	Logger    LoggerConfig `mapstructure:"logger"`
	Jwt       JwtConfig    `mapstructure:"jwt"`
}

var (
	Config = &Conf{
		AppConfig: AppConfig{},
		Mysql:     MysqlConfig{},
		Redis:     RedisConfig{},
		Logger:    LoggerConfig{},
		Jwt:       JwtConfig{},
	}

	once sync.Once

	V = viper.New()
)

func InitConfig(configPath string) {
	once.Do(func() {
		// 设置配置文件
		V.SetConfigFile(configPath)

		// 加载 .yaml 配置
		load()

		V.WatchConfig()
		V.OnConfigChange(func(in fsnotify.Event) {
			load()
		})

		// 检查jwtSecretKey
		checkJwtSecretKey()
	})
}

// 加载配置
func load() {
	if err := V.ReadInConfig(); err != nil {
		fmt.Println("读取配置错误：" + err.Error())
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("未找到配置: " + err.Error())
		} else {
			panic("读取配置出错：" + err.Error())
		}
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := V.Unmarshal(Config); err != nil {
		panic(err)
	}
}

// checkJwtSecretKey 检查jwtSecretKey
func checkJwtSecretKey() {
	// 自动生成JWT secretKey
	if Config.Jwt.SecretKey == "" {
		Config.Jwt.SecretKey = untils.RandString(64)
		V.Set("jwt.secret_key", Config.Jwt.SecretKey)
		err := V.WriteConfig()
		if err != nil {
			panic("自动生成JWT secretKey失败: " + err.Error())
		}
	}
}
