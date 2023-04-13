package conf

import "github.com/spf13/viper"

// 需要监控的 service，用来获取对应的 PID 和 配置文件路径
type Service struct {
	ServiceName       string `mapstructure:"name"`
	ServiceConfigPath string `mapstructure:"config_path"`
}

// 日志配置
type Log struct {
	FileName string `mapstructure:"filename"`
	Max_age  int    `mapstructure:"max_age"`
}

// 指定 exporter 监控端口
type Web struct {
	Addr string `mapstructure:"addr"`
}

type Options struct {
	Service Service `mapstructuer:"service"`
	Web     Web     `mapstructuer:"web"`
	Log     Log     `mapstructuer:"log"`
}

// 解析配置文件并返回解析后的 options 结构体
func ParseConfig(path string) (*Options, error) {
	conf := viper.New()
	conf.SetDefault("web.addr", ":19090")
	conf.SetConfigFile(path)

	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	options := &Options{}
	if err := conf.Unmarshal(options); err != nil {
		return nil, err
	}

	return options, nil
}
