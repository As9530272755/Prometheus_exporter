package etc

import "github.com/spf13/viper"

type Options struct {
	Conf struct {
		Kubeconfig string `mapstructure:"kubeconfig"`
		Port       string `mapstructure:"port"`
	} `mapstructure:"conf"`
}

func ParseConfig(Path string) (*Options, error) {
	conf := viper.New()

	conf.SetDefault("conf.port", ":8080")
	conf.SetConfigFile(Path)
	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	option := &Options{}

	if err := conf.Unmarshal(&option); err != nil {
		return nil, err
	}
	return option, nil
}
