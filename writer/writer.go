package writer

// ZEntry 日志内容
type ZEntry struct {
	Loglevel string // 日志级别
	Message  string // 日志内容
}

type Writer interface {
	// GetLogLevel 获取日志级别
	GetLogLevel() string
	// Check 日志检查，过滤掉不必要的日志输入
	Check(*ZEntry) bool
	// Write 日志输出内容
	Write(message string, attr map[string]interface{}) error
}
