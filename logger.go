package zlog

import (
	"context"
	"fmt"

	"go.elastic.co/apm/module/apmzap"
	"go.uber.org/zap"
)

type Logger struct {
	logger  *zap.Logger
	skipped bool
}

func (log *Logger) Debug(msg string) {
	log.logger.Debug(msg)
}

func (log *Logger) Debugf(template string, args ...interface{}) {
	s := fmt.Sprintf(template, args...)

	log.logger.Debug(s)
}

func (log *Logger) Info(msg string) {
	log.logger.Info(msg)
}

func (log *Logger) Infof(template string, args ...interface{}) {
	s := fmt.Sprintf(template, args...)

	log.logger.Info(s)
}

func (log *Logger) Warn(msg string) {
	log.logger.Warn(msg)
}

func (log *Logger) Warnf(template string, args ...interface{}) {
	s := fmt.Sprintf(template, args...)

	log.logger.Warn(s)
}

func (log *Logger) Error(msg string) {
	log.logger.Error(msg)
}

func (log *Logger) Errorf(template string, args ...interface{}) {
	s := fmt.Sprintf(template, args...)

	log.logger.Error(s)
}

func (log *Logger) With(fields ...zap.Field) *Logger {
	if len(fields) == 0 {
		return log
	}

	l := log.clone()
	l.logger = l.logger.With(fields...)
	if !l.skipped {
		l.skipped = true
		l.logger = l.logger.WithOptions(zap.AddCallerSkip(-1))
	}

	return l
}

func (log *Logger) WithLabel(labels map[string]string) *Logger {
	fields := make([]zap.Field, 0)
	for k, v := range labels {
		fields = append(fields, zap.String(k, v))
	}

	l := log.clone()
	l.logger = l.logger.With(fields...)
	if !l.skipped {
		l.skipped = true
		l.logger = l.logger.WithOptions(zap.AddCallerSkip(-1))
	}
	return l
}

func (log *Logger) WithContext(ctx context.Context) *Logger {
	l := log.clone()
	l.logger = l.logger.With(apmzap.TraceContext(ctx)...)
	if !l.skipped {
		l.skipped = true
		l.logger = l.logger.WithOptions(zap.AddCallerSkip(-1))
	}
	return l
}

func (log *Logger) clone() *Logger {
	copy := *log
	return &copy
}
