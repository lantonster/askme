package log

import (
	"context"
	"os"
	"time"

	"github.com/lantonster/askme/pkg/tracer"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger = zap.NewExample().Sugar()

type Config struct {
	Level      int8   `yaml:"level" mapstructure:"level"`             // 日志级别: [-1, 5] debug, info, warn, error, dpanic, panic, fatal
	FileName   string `yaml:"file_name" mapstructure:"file_name"`     // 日志文件名, 如: ./log/askme.log
	MaxAge     int    `yaml:"max_age" mapstructure:"max_age"`         // 最大保留天数
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size"`       // 最大文件大小，单位 MB
	MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups"` // 最大保留文件数
	Compress   bool   `yaml:"compress" mapstructure:"compress"`       // 是否压缩
}

func SetLogger(config *Config) {
	// 日志格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 日志级别大写并带有颜色
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 日志文件
	file, _ := os.OpenFile(config.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	writer := lumberjack.Logger{
		Filename:   file.Name(),
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
	}

	// 日志输出：文件和控制台
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(&writer), zapcore.Level(config.Level)),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.Level(config.Level)),
	)

	defaultLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
	defaultLogger.Info("logger init success")
}

// WithContext 为日志记录器添加上下文信息。
// 如果上下文中包含跟踪 id，则将其添加到日志中以增强可追踪性。
func WithContext(c context.Context) *zap.SugaredLogger {
	// 如果上下文为空，直接返回日志记录器，不进行任何修改。
	if c == nil {
		return defaultLogger
	}

	// 尝试从上下文中获取跟踪ID。
	// 如果能够成功获取到字符串类型的跟踪ID，则为日志记录器添加"trace_id"字段。
	if traceId, ok := c.Value(tracer.TraceIdKey).(string); ok {
		return defaultLogger.With(tracer.TraceIdKey, traceId)
	}

	// 如果上下文中没有跟踪ID，或者获取失败，则返回原始日志记录器。
	return defaultLogger
}
