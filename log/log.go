package log

import (
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	mlog   = logrus.New()
	mEntry *logrus.Entry
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
		mEntry = mlog.WithField("topics", []string{logConfig.Kafka.Topic})
	}
}

func Entry() *logrus.Entry {
	if mEntry == nil {
		mEntry = mlog.WithField("topics", []string{"none"})
	}
	return mEntry
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
	hook, err := NewKafkaHook(
		[]logrus.Level{logrus.InfoLevel, logrus.DebugLevel, logrus.ErrorLevel},
		&logrus.JSONFormatter{TimestampFormat: "2006-01-2 15:04:05.000"},
		kafkaConfig.Brokers,
	)
	if err != nil {
		panic(fmt.Sprintf("NewKafkaHook failed with err:%v", err))
	}
	log.AddHook(hook)
}
