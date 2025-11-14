package autoload

type DivisionTime struct {
	MaxAge       int `mapstructure:"max_age"`       // 保留旧文件的最大天数，单位天
	RotationTime int `mapstructure:"rotation_time"` // 多长时间切割一次文件，单位小时
}

type DivisionSize struct {
	MaxSize    int  `mapstructure:"max_size"`    // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups int  `mapstructure:"max_backups"` // 保留旧文件的最大个数
	MaxAge     int  `mapstructure:"max_age"`     // 保留旧文件的最大天数
	Compress   bool `mapstructure:"compress"`    // 是否压缩/归档旧文件
}

type LoggerConfig struct {
	DefaultDivision string       `mapstructure:"default_division"`
	Filename        string       `mapstructure:"file_name"`
	DivisionTime    DivisionTime `mapstructure:"division_time"`
	DivisionSize    DivisionSize `mapstructure:"division_size"`
}
