package main

import (
	"github.com/vvfly/zlog"
	"go.uber.org/zap"
)

const (
	debugLogContent = "这是一个调试日志"
	infoLogContent  = "这是一个信息日志"
	warnLogContent  = "这是一个告警日志"
	errLogContent   = "这是一个错误日志"

	corpID      = "123456"
	logTemplate = "corpID=%s, content=%s"
)

func main() {
	// 日志打印
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
}
