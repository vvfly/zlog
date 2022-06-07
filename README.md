zlog是一款底层基于zap的日志记录库,在兼容旧版nlog功能的同时，在使用上更加简单，性能上表现更好。  
zlog专注于日志输出，目前支持日志输出到控制台，文件，sentry，apm，使用者可以自定义日志输出地。  
zlog支持自定义日志检查，默认情况下，zlog会根据日志级别判断日志是否需要输出。

## 一. 支持的日志输出地
- console
- file
- sentry
- APM

## 二. 使用方法
zlog把不同的日志输出地看作是不同的writer，通过注册相应的writer，来控制日志输出。
### 1. console
zlog可以把日志输出到控制台，默认输出到控制台的标准输出。  
可以直接安装zlog使用，也可以通过配置自定义安装使用
#### 1.1 直接安装使用
- 安装zlog  
  `import "github.com/vvfly/zlog"`
- 日志记录  
  zlog.Info(...)  
- 说明：这种情况下，日志级别默认是info

#### 1.2 自定义配置安装使用
- 安装zlog  
  `import "github.com/vvfly/zlog"`
- Console配置项说明
    - LogLevel                string // 日志级别，可选，默认debug
    - CheckFunc func(*ZEntry) bool // 日志检查函数，可选参数
- 通过配置Console配置项初始化zlog
```go
   consoleWriter := writer.NewConsoleWriter(&writer.ConsoleConfig{
        LogLevel: "debug",
    })
   err := zlog.Use(consoleWriter)
```
- 日志记录  
  zlog.Info(...)

### 2. file
zlog可以把日志输出到文件，并通过配置自动对日志文件进行管理，例如日志文件备份，日志文件删除，日志文件大小控制等。
- 安装zlog  
  `import "github.com/vvfly/zlog"`
- File配置项说明
  - Filename    string // 日志文件, 默认"./log/zlog.log",可选
  - MaxFileSize int    // 日志文件单个文件最大大小，单位为MB，默认100,可选
  - MaxBackups  int    // 日志文件最大历史保留份数，默认5,可选
  - MaxAge      int    // 日志文件最长存活时间，单位为天， 默认30,可选
  - Compress    bool   // 历史日志压缩保存，默认false,可选
  - LogLevel    string // 日志级别，默认debug,可选
  - CheckFunc func(*ZEntry) bool // 日志检查函数，可选参数
- 通过配置File配置项初始化zlog
```go
    fileWriter := writer.NewFileWriter(&writer.FileConfig{
        Filename:    "./log/zlog.log",
        MaxFileSize: 100,
        MaxBackups:  5,
        MaxAge:      30,
        Compress:    false,
        LogLevel:    "debug",
    })
    err := zlog.Use(fileWriter)
```
- 日志记录  
  zlog.Info(...)
- 说明  
写本地文件，会提升系统cache缓存使用，请合理配置。例如日志文件大小MaxFileSize配置为100，备份数MaxBackups配置为5，则系统缓存使用会提升500-600M。

### 3. sentry
zlog可以把日志文件上传到sentry，用户可以在sentry ui上查看相应服务的运行情况
- 申请服务的sentry地址(DSN)
- 安装zlog  
    `import "github.com/vvfly/zlog"`
- Sentry配置项说明
  - DSN              string // DSN地址，必填
  - AttachStacktrace bool   // 是否追加堆栈信息,默认否,可选
  - ServerName       string // 服务名， 可选
  - LogLevel         string // 写入Sentry的日志级别,默认info级别，可选
  - CheckFunc        func(*ZEntry) bool // 日志检查，可选参数
- 通过配置Sentry配置项初始化zlog
```go
    sentryWriter := writer.NewSentryWriter(&writer.SentryConfig{
      DSN:              "https://fdec51f7b52f420b98262f3338afdeaf@sentry.ops.weibanzhushou.com/87",
      AttachStacktrace: true,
      ServerName:       "magic-upload",
      LogLevel:         "error",
    })
    err := zlog.Use(sentryWriter)
```
- 日志记录  
  zlog.Error(...)


### 4. APM
zlog可以把日志上传到apm，用户可以在apm ui上查看相应服务的运行情况;  
zlog支持`WithContext()`进行apm链路追踪以及`WithLabel()`来展示APM标签。
- 申请项目的apm资源，导入相应环境变量
- 启动filebeat
- 安装zlog  
  `import "github.com/vvfly/zlog"`
- Apm配置项说明
  - CheckFunc        func(*ZEntry) bool // 日志检查，可选参数
- 通过配置Apm配置项初始化zlog
```go
    apmWriter := writer.NewApmWriter(&writer.ApmConfig{})
    err := zlog.Use(apmWriter)
```
- 日志记录  
  zlog.Info(...)

### 5. 添加结构化内容 Field
zlog提供`With(fields ...zap.Field)`函数来添加结构化内容Field到日志中
```go
zlog.With(zap.String("objectID", "123456789")).Info("success")
```

### 6. 打标签Label
zlog提供`WithLable(map[string]string)`函数进行打标签,标签内容以Filed形式添加到日志中
```go
labels := make(map[string]string)
labels["spanName"] = "errHandler"
labels["spanType"] = "GET"
labels["corpID"] = "12345678"

zlog.WithLabel(labels).Error(errLogContent)
```

### 7. 自定义日志检查
zlog提供自定义日志检查，过滤日志输入;  
默认情况下，zlog会自动检查日志级别, 只有较高级别的日志才会被输出。
```go
func checkFileLog(entry *writer.ZEntry) bool {
  // 日志内容检查
  if strings.Contains(entry.Message, "password") {
    return false
  }
  
  // 日志级别检查
  level, err := zapcore.ParseLevel(entry.Loglevel)
  if err != nil || level < zapcore.ErrorLevel {
    return false
  }
  
  return true
}

zlog.Use(&writer.FileWriter{
    CheckFunc:   checkFileLog,
})
```

### 8. 日志级别
zlog提供`debug info warn error`四种日志级别
- 日志输出
```go
Debug(msg string)
Info(msg string)
Warn(msg string)
Error(msg string)
```
- 自定义格式化日志输出
```go
Debugf(template string, args ...interface{})
Infof(template string, args ...interface{})
Warnf(template string, args ...interface{})
Errorf(template string, args ...interface{})
```
