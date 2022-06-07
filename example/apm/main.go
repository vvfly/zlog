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
	errLogContent = "这是一个错误日志"

	corpID      = "123456"
	logTemplate = "corpID=%s, content=%s"
)

// registeApmWriter 注册apm writer
func registeApmWriter() {
	// apm
	zlog.Use(writer.NewApmWriter(&writer.ApmConfig{}))
}

func main() {

	// init log
	registeApmWriter()

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
	// error日志输出
	zlog.Error(errLogContent)
	zlog.Errorf(logTemplate, corpID, errLogContent)

	// 添加结构化内容Field
	zlog.With(zap.String("name", "zhangsan"), zap.Int("age", 28)).Error(errLogContent)

	// 添加标签
	labels := make(map[string]string)
	labels["spanName"] = "zlogHandler"
	labels["spanType"] = "GET"
	labels["corpID"] = corpID
	zlog.WithLabel(labels).Error(errLogContent)

	c.String(http.StatusOK, "success")
}
