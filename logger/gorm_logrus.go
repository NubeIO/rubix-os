package logger

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"time"

	log "github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
)

type logger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	LogLevel              gormLogger.LogLevel
}

func New() *logger {
	return &logger{
		SkipErrRecordNotFound: true,
	}
}

func (l *logger) SetLogMode(logLevel string) gormLogger.Interface {
	var level gormLogger.LogLevel
	switch logLevel {
	case "INFO":
		level = gormLogger.Info
	case "WARN":
		level = gormLogger.Warn
	case "ERROR":
		level = gormLogger.Error
	default:
		level = gormLogger.Silent
	}
	l.LogLevel = level
	return l
}

func (l *logger) LogMode(_ gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *logger) Info(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Infof(s, args)
}

func (l *logger) Warn(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Warnf(s, args)
}

func (l *logger) Error(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Errorf(s, args)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	sql, _ := fc()
	fields := log.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[log.ErrorKey] = err
		log.WithContext(ctx).WithFields(fields).Errorf("%s (%s)", sql, elapsed)
		return
	}
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		log.WithContext(ctx).WithFields(fields).Warnf("%s (%s)", sql, elapsed)
		return
	}
	if l.LogLevel == gormLogger.Info {
		log.WithContext(ctx).WithFields(fields).Infof("%s (%s)", sql, elapsed)
	}
}
