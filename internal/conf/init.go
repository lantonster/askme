package conf

import (
	"strings"

	"github.com/google/wire"
	"github.com/lantonster/askme/configs"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var ProviderSetConfig = wire.NewSet(NewConfig)

func NewConfig() *Config {
	config := initDefault()

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}

func initDefault() *Config {
	config := &Config{}
	yaml.Unmarshal(configs.Config, config)
	return config
}
