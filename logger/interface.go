package logger

import (
	"context"
	"errors"
	"time"
)

// 常用错误定义
var (
	ErrRecordNotFound = errors.New("record not found")
)

// LogLevel 日志级别
type LogLevel int

const (
	// Silent 静默模式，不输出任何日志
	Silent LogLevel = iota + 1
	// Error 只输出错误日志
	Error
	// Warn 输出警告和错误日志
	Warn
	// Info 输出所有日志（包括 SQL 语句）
	Info
)

// Interface 日志接口，参考 GORM 的设计
type Interface interface {
	// LogMode 设置日志级别
	LogMode(level LogLevel) Interface

	// Info 输出信息日志
	Info(ctx context.Context, msg string, data ...interface{})

	// Warn 输出警告日志
	Warn(ctx context.Context, msg string, data ...interface{})

	// Error 输出错误日志
	Error(ctx context.Context, msg string, data ...interface{})

	// Trace 追踪 SQL 执行
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

// Config 日志配置
type Config struct {
	// SlowThreshold 慢查询阈值
	SlowThreshold time.Duration
	// Colorful 是否启用彩色输出
	Colorful bool
	// IgnoreRecordNotFoundError 是否忽略记录未找到错误
	IgnoreRecordNotFoundError bool
	// ParameterizedQueries 是否显示参数化查询（隐藏参数）
	ParameterizedQueries bool
	// LogLevel 日志级别
	LogLevel LogLevel
}

// Writer 日志写入器接口
type Writer interface {
	Printf(string, ...interface{})
}
