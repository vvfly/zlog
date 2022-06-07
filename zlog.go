package zlog

import (
	"context"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/vvfly/zlog/writer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultKeyGoVersion       = "go.version"
	defaultKeyHostName        = "hostname"
	defaultKeyBuildAppVersion = "build.version"
	defaultKeyBuildUser       = "build.user"
	defaultKeyBuildHost       = "build.host"
	defaultKeyBuildTime       = "build.time"
	defaultKeyApplicationName = "application.name"
)

var (
	goVersion       string
	hostName        string
	buildAppVersion string
	buildUser       string
	buildHost       string
	buildTime       string
)

// 全局log对象，可以使用Use()方法重置
var _logger *Logger

func init() {
	// _logger init
	_logger = &Logger{
		logger: initZapLogger("info"),
	}

	// env
	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	goVersion = runtime.Version()
}

// initZapLogger create a zap logger
func initZapLogger(level string) *zap.Logger {
	// writer: 默认标准输出
	zwriter := zapcore.AddSync(os.Stdout)

	// encoder: 输出JSON格式
	encoder := getJsonEncoder()

	core := zapcore.NewCore(encoder, zwriter, writer.GetLoggerLevel(level))

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
}

// Use 注册writer
func Use(writers ...writer.Writer) error {

	zCore := &zCore{
		level:   writer.GetLowLogLevel(writers),
		enc:     getJsonEncoder(),
		writers: writers,
		attr:    make(map[string]interface{}),
	}
	logger := zap.New(zCore, zap.AddCaller(), zap.AddCallerSkip(2))

	// with default fields
	logger = logger.With(defaultFields()...)

	_logger = &Logger{
		logger: logger,
	}
	return nil
}

func Debug(msg string) {
	_logger.Debug(msg)
}

func Debugf(template string, args ...interface{}) {
	_logger.Debugf(template, args...)
}

func Info(msg string) {
	_logger.Info(msg)
}

func Infof(template string, args ...interface{}) {
	_logger.Infof(template, args...)
}

func Warn(msg string) {
	_logger.Warn(msg)
}

func Warnf(template string, args ...interface{}) {
	_logger.Warnf(template, args...)
}

func Error(msg string) {
	_logger.Error(msg)
}

func Errorf(template string, args ...interface{}) {
	_logger.Errorf(template, args...)
}

func With(fields ...zap.Field) *Logger {
	return _logger.With(fields...)
}

func WithLabel(labels map[string]string) *Logger {
	return _logger.WithLabel(labels)
}

func WithContext(ctx context.Context) *Logger {
	return _logger.WithContext(ctx)
}

func getJsonEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2006-01-02 15:04:05.000000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func defaultFields() []zap.Field {
	var fields []zap.Field
	projectName := os.Getenv("ELASTIC_APM_SERVICE_NAME")
	if len(projectName) == 0 {
		projectName = os.Args[0]
		projectName = strings.Trim(projectName, "./")
	}
	if len(projectName) == 0 {
		projectName = "empty-project-name"
	}
	fields = append(fields, zap.String(defaultKeyApplicationName, projectName))
	if buildAppVersion != "" {
		fields = append(fields, zap.String(defaultKeyBuildAppVersion, buildAppVersion))
	}
	if buildTime != "" {
		fields = append(fields, zap.String(defaultKeyBuildTime, buildTime))
	}
	if buildHost != "" {
		fields = append(fields, zap.String(defaultKeyBuildHost, buildHost))
	}
	if buildUser != "" {
		fields = append(fields, zap.String(defaultKeyBuildUser, buildUser))
	}
	if hostName != "" {
		fields = append(fields, zap.String(defaultKeyHostName, hostName))
	}
	if goVersion != "" {
		fields = append(fields, zap.String(defaultKeyGoVersion, goVersion))
	}
	return fields
}
