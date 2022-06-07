package main

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vvfly/zlog"
	"github.com/vvfly/zlog/writer"
)

const (
	debugLogContent = "这是一个调试日志"
	infoLogContent  = "这是一个信息日志"
	warnLogContent  = "这是一个告警日志"
	errLogContent   = "这是一个错误日志,包含password敏感信息"

	corpID      = "123456"
	logTemplate = "corpID=%s, content=%s"
)

// checkFileLog 自定义日志检查, 如果日志内容中有password敏感信息，就跳过该条日志输出
func checkFileLog(entry *writer.ZEntry) bool {
	// 日志内容检查（过滤敏感信息）
	if strings.Contains(entry.Message, "password") {
		return false
	}

	// 日志级别检查
	level, err := zapcore.ParseLevel(entry.Loglevel)
	if err != nil || level <= zapcore.DebugLevel { // 小于等于debug级别日志会被过滤
		return false
	}

	return true
}

// registeLogWriter 注册writer
func registeLogWriter() {
	// file
	fileWriter := writer.NewFileWriter(&writer.FileConfig{
		Filename:    "./log/zlog.log",
		MaxFileSize: 100,
		MaxBackups:  5,
		MaxAge:      30,
		Compress:    false,
		LogLevel:    "debug",
		CheckFunc:   checkFileLog,
	})

	zlog.Use(fileWriter)
}

func main() {

	// init log
	registeLogWriter()

	logPrint()
}

func logPrint() {
	// debug日志输出
	zlog.Debug(debugLogContent)
	zlog.Debugf(logTemplate, corpID, debugLogContent)

	// info日志输出
	zlog.Info(infoLogContent)
	zlog.Infof(logTemplate, corpID, infoLogContent)

	// warn日志输出
	zlog.Warn(warnLogContent)
	zlog.Warnf(logTemplate, corpID, warnLogContent)

	// error日志输出
	zlog.Error(errLogContent)
	zlog.Errorf(logTemplate, corpID, errLogContent)

	// 添加结构化内容Field
	zlog.With(zap.String("name", "zhangsan"), zap.Int("age", 28)).Info(infoLogContent)

	// 添加标签
	labels := make(map[string]string)
	labels["spanName"] = "zlogHandler"
	labels["spanType"] = "GET"
	labels["corpID"] = corpID
	zlog.WithLabel(labels).Info(infoLogContent)
}
