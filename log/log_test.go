package log

import (
	"testing"
	"time"
)

func TestLocalSystemLogger(T *testing.T) {
	config := LogConfig{
		Kafka: nil,
		Path:  "logs",
		Level: "debug",
	}
	InitLog(config)
	l := Entry()
	l.Info("this is info")
	l.Debug("this is debug")
	l.Error("this is error")
}

func TestKafkaLogger(T *testing.T) {
	config := LogConfig{
		Kafka: &KafkaConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "nodelog",
		},
		Path:  "logs",
		Level: "debug",
	}
	InitLog(config)
	l := Entry()
	l.Info("this is info")
	l.Debug("this is debug")
	l.Error("this is error")
	time.Sleep(time.Second)
}
