package log

import (
	"testing"
)

func TestLocalSystemLogger(T *testing.T) {
	config := LogConfig{
		Kafka: nil,
		Path:  "logs",
		Level: "debug",
	}
	InitLog(config)
	Info("this is info")
	Debug("this is debug")
	Error("this is error")
}

func TestKafkaLogger(T *testing.T) {
	config := LogConfig{
		Kafka: nil,
		Path:  "logs",
		Level: "debug",
	}
	InitLog(config)
	Info("this is info")
	Debug("this is debug")
	Error("this is error")
}
