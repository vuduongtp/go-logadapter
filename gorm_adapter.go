package logadapter

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger model
type GormLogger struct {
	*Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

// NewGormLogger new gorm logger
func NewGormLogger() *GormLogger {
	return &GormLogger{
		SkipErrRecordNotFound: true,
		Debug:                 true,
		Logger:                l,
		SlowThreshold:         time.Second,
		SourceField:           DefaultSourceField,
	}
}

// LogMode get log mode, should set to Silent level in production
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	if level == gormlogger.Silent {
		l.Debug = false
	} else {
		l.Debug = true
	}

	return l
}

// Info log infor
func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Infof(s, args...)
}

// Warn log warn
func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Warnf(s, args...)
}

// Error log error
func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Errorf(s, args...)
}

// Trace log sql trace
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, row := fc()
	trace := logrus.Fields{
		"type":       LogTypeSQL,
		"row":        row,
		"latency_ms": elapsed.Milliseconds(),
		"latency":    elapsed.String(),
	}

	if l.Debug {
		trace["query"] = sql
	}

	fields := mergeLogFields(trace, GetLogFieldFromContext(ctx))

	if l.SourceField != "" {
		fields[l.SourceField] = l.getCaller()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.Logger.WithContext(ctx).WithFields(fields).Error()
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.WithContext(ctx).WithFields(fields).Warn()
		return
	}

	if l.Debug {
		l.Logger.WithContext(ctx).WithFields(fields).Debug()
	}

}
