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

// GormLogAdapter model
type GormLogAdapter struct {
	*Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

// NewGormLogAdapter gorm logrus GormLogAdapter
func NewGormLogAdapter(log *Logger) *GormLogAdapter {
	return &GormLogAdapter{
		SkipErrRecordNotFound: true,
		Debug:                 true,
		Logger:                log,
	}
}

// LogMode function
func (l *GormLogAdapter) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info function
func (l *GormLogAdapter) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Infof(s, args)
}

// Warn function
func (l *GormLogAdapter) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Warnf(s, args)
}

// Error function
func (l *GormLogAdapter) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Errorf(s, args)
}

// Trace function
func (l *GormLogAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
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
