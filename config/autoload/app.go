package autoload

type AppConfig struct {
	AppEnv         string  `mapstructure:"app_env"`
	Debug          bool    `mapstructure:"debug"`
	Language       string  `mapstructure:"language"`
	WatchConfig    bool    `mapstructure:"watch_config"`
	StaticBasePath string  `mapstructure:"base_path"`
	Timezone       *string `mapstructure:"timezone"`
	Port           int     `mapstructure:"port"`
}
