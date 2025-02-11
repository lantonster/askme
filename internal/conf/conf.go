package conf

import (
	"time"

	"github.com/lantonster/askme/pkg/i18n"
	"github.com/lantonster/askme/pkg/log"
)

type Config struct {
	Server  *Server      `yaml:"server" mapstructure:"server"`   // 服务
	Swagger *Swagger     `yaml:"swagger" mapstructure:"swagger"` // swagger 配置
	UI      *UI          `yaml:"ui" mapstructure:"ui"`           // ui 配置
	Uploads *Uploads     `yaml:"uploads" mapstructure:"uploads"` // 上传配置
	Logger  *log.Config  `yaml:"logger" mapstructure:"logger"`   // 日志配置
	I18n    *i18n.Config `yaml:"i18n" mapstructure:"i18n"`       // 国际化配置
}

type Server struct {
	Http *Http `yaml:"http" mapstructure:"http"` // http 服务
}

type Http struct {
	Debug           bool          `yaml:"debug" mapstructure:"debug"`                       // 是否开启 debug 模式
	Addr            string        `yaml:"addr" mapstructure:"addr"`                         // 监听地址
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"` // 关闭超时时间
}

type Swagger struct {
	Show     bool   `yaml:"show" mapstructure:"show"`         // 是否开启 swagger
	Protocal string `yaml:"protocal" mapstructure:"protocal"` // 协议: http, https
	Host     string `yaml:"host" mapstructure:"host"`         // 主机
}

type UI struct {
	BaseUrl string `yaml:"base_url" mapstructure:"base_url"` // 基础 url
}
