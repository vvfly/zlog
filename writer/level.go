package writer

import (
	"go.uber.org/zap/zapcore"
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

func checkLogLevel(entryLv, checkLv string) bool {
	return GetLoggerLevel(entryLv) >= GetLoggerLevel(checkLv)
}

func GetLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

// GetLowLogLevel 获取writer中最低日志级别
func GetLowLogLevel(writers []Writer) zapcore.Level {
	var level = zapcore.ErrorLevel
	for i := range writers {
		logLevel := writers[i].GetLogLevel()
		zapLogLevel := GetLoggerLevel(logLevel)
		if zapLogLevel <= level {
			level = zapLogLevel
		}
	}
	return level
}
