package conf

import "time"

type Config struct {
	Server *Server `yaml:"server" mapstructure:"server"` // 服务
	UI     *UI     `yaml:"ui" mapstructure:"ui"`         // ui 配置
}

type Server struct {
	Http *Http `yaml:"http" mapstructure:"http"` // http 服务
}

type Http struct {
	Debug           bool          `yaml:"debug" mapstructure:"debug"`                       // 是否开启 debug 模式
	Addr            string        `yaml:"addr" mapstructure:"addr"`                         // 监听地址
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"` // 关闭超时时间
}

type UI struct {
	BaseUrl string `yaml:"base_url" mapstructure:"base_url"` // 基础 url
}
