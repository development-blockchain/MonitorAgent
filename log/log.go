package log

import (
	"context"
	"fmt"
	"github.com/develope/MonitorAgent/common"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	mlog = logrus.New()
)

type KafkaConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type LogConfig struct {
	Kafka *KafkaConfig
	Path  string `json:"path"`  // local logs file store path
	Level string `json:"level"` // log level
}

func logLevel(level string) logrus.Level {
	switch level {
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func InitLog(logConfig LogConfig) {
	// standard setting
	mlog.SetLevel(logLevel(logConfig.Level))
	mlog.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-2 15:04:05.000"})

	// file system logger setting
	localFilesystemLogger(mlog, logConfig.Path)
	if logConfig.Kafka != nil {
		// kafka logger setting
		kafkaLogger(mlog, logConfig.Kafka)
	}
}

func Logger() *logrus.Logger {
	return mlog
}

func logWriter(logPath string) *rotatelogs.RotateLogs {
	logFullPath := logPath
	logwriter, err := rotatelogs.New(
		logFullPath+".%Y%m%d",
		rotatelogs.WithLinkName(logFullPath),
		rotatelogs.WithRotationSize(100*1024*1024), // 100MB
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return logwriter
}

func localFilesystemLogger(log *logrus.Logger, logPath string) {
	lfHook := lfshook.NewHook(logWriter(logPath), &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-2 15:04:05.000"})
	log.AddHook(lfHook)
}

func kafkaLogger(log *logrus.Logger, kafkaConfig *KafkaConfig) {
	var nodeid, _ = os.Hostname()
	if nodeip, err := common.GetExternal(); err == nil {
		nodeid = fmt.Sprintf("node_%s", strings.Replace(nodeip, ".", "_", -1))
	}
	hook, err := NewKafkaHook(
		nodeid,
		[]logrus.Level{logrus.InfoLevel, logrus.DebugLevel, logrus.ErrorLevel},
		&logrus.JSONFormatter{},
		kafkaConfig.Brokers,
	)
	if err != nil {
		panic(fmt.Sprintf("NewKafkaHook failed with err:%v", err))
	}
	log.AddHook(hook)
}

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func WithField(key string, value interface{}) *logrus.Entry {
	return mlog.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return mlog.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func WithError(err error) *logrus.Entry {
	return mlog.WithError(err)
}

// Add a context to the log entry.
func WithContext(ctx context.Context) *logrus.Entry {
	return mlog.WithContext(ctx)
}

// Overrides the time of the log entry.
func WithTime(t time.Time) *logrus.Entry {
	return mlog.WithTime(t)
}

func Logf(level logrus.Level, format string, args ...interface{}) {
	mlog.Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) {
	mlog.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	mlog.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	mlog.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	mlog.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	mlog.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	mlog.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	mlog.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	mlog.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	mlog.Panicf(format, args...)
}

func Log(level logrus.Level, args ...interface{}) {
	mlog.Log(level, args...)
}

func LogFn(level logrus.Level, fn logrus.LogFunction) {
	mlog.LogFn(level, fn)
}

func Trace(args ...interface{}) {
	mlog.Trace(args...)
}

func Debug(args ...interface{}) {
	mlog.Debug(args...)
}

func Info(args ...interface{}) {
	mlog.Info(args...)
}

func Print(args ...interface{}) {
	mlog.Print(args...)
}

func Warn(args ...interface{}) {
	mlog.Warn(args...)
}

func Warning(args ...interface{}) {
	mlog.Warning(args...)
}

func Error(args ...interface{}) {
	mlog.Error(args...)
}

func Fatal(args ...interface{}) {
	mlog.Fatal(args...)
}

func Panic(args ...interface{}) {
	mlog.Panic(args...)
}

func TraceFn(fn logrus.LogFunction) {
	mlog.TraceFn(fn)
}

func DebugFn(fn logrus.LogFunction) {
	mlog.DebugFn(fn)
}

func InfoFn(fn logrus.LogFunction) {
	mlog.InfoFn(fn)
}

func PrintFn(fn logrus.LogFunction) {
	mlog.PrintFn(fn)
}

func WarnFn(fn logrus.LogFunction) {
	mlog.WarnFn(fn)
}

func WarningFn(fn logrus.LogFunction) {
	mlog.WarningFn(fn)
}

func ErrorFn(fn logrus.LogFunction) {
	mlog.ErrorFn(fn)
}

func FatalFn(fn logrus.LogFunction) {
	mlog.FatalFn(fn)
}

func PanicFn(fn logrus.LogFunction) {
	mlog.PanicFn(fn)
}

func Logln(level logrus.Level, args ...interface{}) {
	mlog.Logln(level, args...)
}

func Traceln(args ...interface{}) {
	mlog.Traceln(args...)
}

func Debugln(args ...interface{}) {
	mlog.Debugln(args...)
}

func Infoln(args ...interface{}) {
	mlog.Infoln(args...)
}

func Println(args ...interface{}) {
	mlog.Println(args...)
}

func Warnln(args ...interface{}) {
	mlog.Warnln(args...)
}

func Warningln(args ...interface{}) {
	mlog.Warningln(args...)
}

func Errorln(args ...interface{}) {
	mlog.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	mlog.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	mlog.Panicln(args...)
}
