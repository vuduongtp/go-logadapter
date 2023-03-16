package logadapter

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
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
		SourceField:           DefaultGormSourceField,
	}
}

// LogMode get log mode
func (l *GormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info log infor
func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Infof(s, args)
}

// Warn log warn
func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Warnf(s, args)
}

// Error log error
func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Errorf(s, args)
}

// Trace log sql trace
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, row := fc()
	fields := logrus.Fields{}
	fields["type"] = LogTypeSQL
	fields["row"] = row
	fields["latency_ms"] = elapsed.Milliseconds()
	fields["query"] = sql

	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.Logger.WithContext(ctx).WithFields(fields).Errorf(`[row:%d] [%s] [%s]`, row, elapsed, sql)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.WithContext(ctx).WithFields(fields).Warnf(`[row:%d] [%s] [%s]`, row, elapsed, sql)
		return
	}

	if l.Debug {
		l.Logger.WithContext(ctx).WithFields(fields).Debugf(`[row:%d] [%s] [%s]`, row, elapsed, sql)
	}
}
