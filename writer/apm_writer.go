package writer

import (
	"go.elastic.co/apm"
)

const (
	// FieldKeyTraceID is the field key for the trace ID.
	FieldKeyTraceID = "trace.id"

	// FieldKeyTransactionID is the field key for the transaction ID.
	FieldKeyTransactionID = "transaction.id"

	// FieldKeySpanID is the field key for the span ID.
	FieldKeySpanID = "span.id"
)

// ApmConfig 用户自定义apm配置
type ApmConfig struct {
	CheckFunc func(*ZEntry) bool
}

// apmWriter 实现writer接口
type apmWriter struct {
	// Tracer is the apm.Tracer to use for reporting errors.
	Tracer *apm.Tracer
	// 日志级别，默认error
	LogLevel string
	opt      *ApmConfig
}

func (a *apmWriter) GetLogLevel() string {
	return a.LogLevel
}

func (a *apmWriter) Check(entry *ZEntry) bool {
	if a.opt.CheckFunc != nil {
		return a.opt.CheckFunc(entry)
	}

	return checkLogLevel(entry.Loglevel, a.LogLevel)
}

func (a *apmWriter) Write(message string, attr map[string]interface{}) error {

	tracer := a.Tracer
	errlog := tracer.NewErrorLog(apm.ErrorLogRecord{
		Message: message,
	})
	errlog.Handled = true
	errlog.SetStacktrace(1)

	// set attr
	for k, v := range attr {
		switch k {
		case FieldKeyTraceID:
			errlog.TraceID = v.(apm.TraceID)
			continue
		case FieldKeyTransactionID:
			errlog.TransactionID = v.(apm.SpanID)
			continue
		case FieldKeySpanID:
			spanID := v.(apm.SpanID)
			if spanID.Validate() == nil {
				errlog.ParentID = spanID
			}
			continue
		default:
			errlog.Context.SetLabel(k, v)
		}
	}

	errlog.Send()

	return nil
}

func NewApmWriter(opt *ApmConfig) Writer {
	apmWriter := apmWriter{
		Tracer:   apm.DefaultTracer,
		LogLevel: "error", // apm agent 默认只error级别及以上日志
		opt:      opt,
	}
	return &apmWriter
}
