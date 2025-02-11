package conf

import (
	"context"
	"strings"

	"github.com/google/wire"
	"github.com/lantonster/askme/configs"
	"github.com/lantonster/askme/pkg/i18n"
	"github.com/lantonster/askme/pkg/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var ProviderSet = wire.NewSet(NewConfig)

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

	// 初始化日志
	log.SetLogger(config.Logger)
	log.WithContext(context.Background()).Info("config init success")

	// 初始化国际化
	i18n.SetTranslator(config.I18n)
	log.WithContext(context.Background()).Info("i18n init success")

	return config
}

func initDefault() *Config {
	config := &Config{}
	yaml.Unmarshal(configs.Config, config)
	return config
}
