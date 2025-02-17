package log

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormLogger(debug bool) *GormLogger {
	return &GormLogger{
		Logger:        defaultLogger,
		SlowThreshold: 200 * time.Millisecond,
		Debug:         debug,
	}
}

type GormLogger struct {
	Logger        Logger
	SlowThreshold time.Duration
	Debug         bool
}

// LogMode 实现 logger.Interface 的 LogMode 方法
func (l GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return GormLogger{
		Logger:        l.Logger,
		SlowThreshold: l.SlowThreshold,
		Debug:         l.Debug,
	}
}

// Info 实现 logger.Interface 的 Info 方法
func (l GormLogger) Info(c context.Context, str string, args ...interface{}) {
	l.logger(c).Debugf(str, args...)
}

// Warn 实现 logger.Interface 的 Warn 方法
func (l GormLogger) Warn(c context.Context, str string, args ...interface{}) {
	l.logger(c).Warnf(str, args...)
}

// Error 实现 logger.Interface 的 Error 方法
func (l GormLogger) Error(c context.Context, str string, args ...interface{}) {
	l.logger(c).Errorf(str, args...)
}

// Trace 实现 logger.Interface 的 Trace 方法
func (l GormLogger) Trace(c context.Context, begin time.Time, fc func() (string, int64), err error) {

	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 请求和返回条数
	sql, rows := fc()

	// 携带基础信息
	logger := l.logger(c).With("sql", sql, "rows", rows, "time", fmt.Sprintf("%v", time.Since(begin)))

	// Gorm 错误
	if err != nil {
		// 记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Database ErrRecordNotFound")
		} else {
			logger.Error("Database Error: ", err)
		}
		return
	}

	// 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logger.Warn("Database Slow Log")
		return
	}

	// 测试环境记录所有 SQL 请求
	if l.Debug {
		logger.Debug("Database Query")
	}
}

// logger 内用的辅助方法，确保 Zap 内置信息 Caller 的准确性（如 paginator/paginator.go:148）
func (l GormLogger) logger(c context.Context) *zap.SugaredLogger {
	// 跳过 gorm 内置的调用
	var (
		gormPackage    = filepath.Join("gorm.io", "gorm")
		zapgormPackage = filepath.Join("moul.io", "zapgorm2")
	)

	// 减去一次封装，以及一次在 logger 初始化里添加 zap.AddCallerSkip(1)
	clone := l.Logger.WithContext(c).Desugar().WithOptions(zap.AddCallerSkip(-2))

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			// 返回一个附带跳过行号的新的 zap logger
			return clone.WithOptions(zap.AddCallerSkip(i)).Sugar()
		}
	}

	return l.Logger.WithContext(c)
}
