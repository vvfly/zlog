package writer

import (
	"strings"

	"github.com/getsentry/sentry-go"
)

// 用户自定义sentry配置
type SentryConfig struct {
	DSN              string             // DSN地址，必填
	AttachStacktrace bool               // 是否追加堆栈信息,默认否,可选
	ServerName       string             // 服务名， 可选
	LogLevel         string             // 写入Sentry的日志级别,默认info级别，可选
	CheckFunc        func(*ZEntry) bool // 日志检查，可选参数
}

// sentryWriter实现writer接口
type sentryWriter struct {
	opt *SentryConfig
}

func (s *sentryWriter) GetLogLevel() string {
	return s.opt.LogLevel
}

func (s *sentryWriter) Check(entry *ZEntry) bool {
	if s.opt.CheckFunc != nil {
		return s.opt.CheckFunc(entry)
	}

	return checkLogLevel(entry.Loglevel, s.opt.LogLevel)
}

func (s *sentryWriter) Write(message string, attr map[string]interface{}) error {
	sentry.CaptureMessage(message)
	return nil
}

func NewSentryWriter(opt *SentryConfig) Writer {
	// LogLevel预处理
	opt.LogLevel = strings.ToLower(strings.TrimSpace(opt.LogLevel))

	// 初始化 sentry client
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              opt.DSN,
		AttachStacktrace: opt.AttachStacktrace,
		ServerName:       opt.ServerName,
	})
	if err != nil {
		panic(err)
	}

	sentryWriter := sentryWriter{
		opt: opt,
	}

	return &sentryWriter
}
