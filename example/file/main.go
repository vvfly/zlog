package main

import (
	"go.uber.org/zap"

	"github.com/vvfly/zlog"
	"github.com/vvfly/zlog/writer"
)

const (
	debugLogContent = "这是一个调试日志"
	infoLogContent  = "这是一个信息日志"
	warnLogContent  = "这是一个告警日志"
	errLogContent   = "这是一个错误日志"

	corpID      = "123456"
	logTemplate = "corpID=%s, content=%s"
)

// registeFileWriter 注册file writer
func registeFileWriter() {
	// file
	fileWriter := writer.NewFileWriter(&writer.FileConfig{
		Filename:    "./log/zlog.log",
		MaxFileSize: 100,
		MaxBackups:  5,
		MaxAge:      30,
		Compress:    false,
		LogLevel:    "info",
	})

	zlog.Use(fileWriter)
}

func main() {
	// init log
	registeFileWriter()

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
