package config

import (
	"github.com/spf13/viper"
)

// 定义 MySQL 结构体用户 options 解析
type MySql struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

// 定义 log 结构体
type Log struct {
	FileName    string `mapstructure:"filename"`
	Level       string `mapstructure:"level"`
	Max_age     int    `mapstructure:"max_age"`
	Max_size    int    `mapstructure:"max_size"`
	Max_backups int    `mapstructure:"max_backups"`
	Compress    bool   `mapstructure:"compress"`
}

// 定义 web 结构体
type Web struct {
	Addr string `mapstructuer:"addr"`
	Auth struct {
		UserName string `mapstructuer:"username"`
		Password string `mapstructuer:"password"`
	} `mapstructuer:"auth"`
}

// 定义解析的配置文件对象
type Options struct {
	MySql MySql `mapstructuer:"mysql"`
	Web   Web   `mapstructuer:"web"`
	Log   Log   `mapstructuer:"log"`
}

// 解析配置文件并返回解析后的 options 结构体
func ParseConfig(path string) (*Options, error) {
	// 定义一个新对象
	conf := viper.New()

	// 设置 wbe 地址默认值
	conf.SetDefault("web.addr", ":9090")
	conf.SetDefault("mysql.port", "3306")

	// 解析传入的配置文件
	conf.SetConfigFile(path)

	// 解析配置文件
	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	// 将配置文件内容读取到结构体中
	option := &Options{}
	if err := conf.Unmarshal(&option); err != nil {
		return nil, err
	}
	return option, nil
}
