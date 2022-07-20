package ex_config

import (
	"github.com/spf13/viper"
)

type KubeConfig struct {
	ConfigPath string `mapstructure:"path"`
}

type Web struct {
	ListenPort string `mapstructure:"port"`
}

type Log struct {
	FileName string `mapstructure:"filename"`
	Max_age  int    `mapstructure:"max_age"`
	Max_size int    `mapstructure:"max_size"`
	Compress bool   `mapstructure:"compress"`
	Level    string `mapstructure:"level"`
}

type Optins struct {
	KubeConfig KubeConfig `mapstructure:"kubeconfig"`
	Log        Log        `mapstructure:"log"`
	Web        Web        `mapstructure:"web"`
}

func ParseConfig(path string) (*Optins, error) {
	conf := viper.New()

	conf.SetConfigFile(path)

	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	optins := &Optins{}
	if err := conf.Unmarshal(&optins); err != nil {
		return nil, err
	}
	return optins, nil
}
