package zlog

import (
	"github.com/vvfly/zlog/writer"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

// zCore 实现zapcore.Core接口
type zCore struct {
	level   zapcore.Level
	enc     zapcore.Encoder
	writers []writer.Writer
	attr    map[string]interface{}
}

func (c *zCore) Enabled(level zapcore.Level) bool {
	return level >= c.level
}

func (c *zCore) With(fields []zapcore.Field) zapcore.Core {
	newCore := &zCore{
		level:   c.level,
		enc:     c.enc.Clone(),
		writers: c.writers,
		attr:    make(map[string]interface{}),
	}

	// add fields to encoder
	for i := range fields {
		fields[i].AddTo(newCore.enc)
	}

	// add fields to newCore attr
	newCore.fields(fields)

	return newCore
}

func (c *zCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *zCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	var err error

	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}

	message := buf.String()
	buf.Free()

	zEntry := writer.ZEntry{
		Loglevel: ent.Level.String(),
		Message:  message,
	}

	for i := range c.writers {
		if c.writers[i].Check(&zEntry) {
			err = multierr.Append(err, c.writers[i].Write(message, c.attr))
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *zCore) Sync() error {
	return nil
}

func (c *zCore) fields(fields []zapcore.Field) {
	for _, field := range fields {
		switch field.Key {
		case defaultKeyGoVersion,
			defaultKeyHostName,
			defaultKeyBuildAppVersion,
			defaultKeyBuildUser,
			defaultKeyBuildHost,
			defaultKeyBuildTime,
			defaultKeyApplicationName:
			continue
		default:
			if field.Integer > 0 {
				c.attr[field.Key] = field.Integer
			} else if field.String != "" {
				c.attr[field.Key] = field.String
			} else {
				c.attr[field.Key] = field.Interface
			}
		}
	}
}
