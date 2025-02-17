package conf

type MySQL struct {
	Debug    bool   `yaml:"debug" mapstructure:"debug"`       // 是否开启 debug 模式
	Addr     string `yaml:"addr" mapstructure:"addr"`         // host:port
	User     string `yaml:"user" mapstructure:"user"`         // 用户名
	Password string `yaml:"password" mapstructure:"password"` // 密码
	DBName   string `yaml:"db_name" mapstructure:"db_name"`   // 数据库名
}

type Redis struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`         // 地址
	User     string `yaml:"user" mapstructure:"user"`         // 用户名
	Password string `yaml:"password" mapstructure:"password"` // 密码
}
