package main

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/vvfly/zlog"
	"github.com/vvfly/zlog/writer"
	"go.elastic.co/apm/module/apmgin"
)

const (
	debugLogContent = "这是一个调试日志"
	infoLogContent  = "这是一个信息日志"
	warnLogContent  = "这是一个告警日志"
	errLogContent   = "这是一个错误日志"

	corpID      = "123456"
	logTemplate = "corpID=%s, content=%s"
)

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
	})

	// sentry
	sentryWriter := writer.NewSentryWriter(&writer.SentryConfig{
		DSN:              "https://fdec51f7b52f420b98262f3338afdeaf@sentry.ops.weibanzhushou.com/87",
		AttachStacktrace: true,
		ServerName:       "magic-upload",
		LogLevel:         "error",
	})

	// apm
	apmWriter := writer.NewApmWriter(&writer.ApmConfig{})

	zlog.Use(fileWriter, sentryWriter, apmWriter)
}

func main() {

	// init log
	registeLogWriter()

	// gin http server
	r := gin.New()

	// add apm middleware handler for tracing requests and reporting errors.
	r.Use(apmgin.Middleware(r))

	r.GET("/zlog", zlogHandler)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		zlog.Errorf("服务启动错误:", err)
		return
	}
}

func zlogHandler(c *gin.Context) {
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

	c.String(http.StatusOK, "success")
}