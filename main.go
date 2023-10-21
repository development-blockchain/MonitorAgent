package main

import (
	"github.com/develope/MonitorAgent/log"
	"github.com/develope/MonitorAgent/metrics"
)

func main() {
	metrics.StartMetrics(9000)
	config := log.LogConfig{
		//Kafka: &log.KafkaConfig{
		//	Brokers: []string{"localhost:9092"},
		//	Topic:   "topic_1",
		//},
		Path:  "logs",
		Level: "debug",
	}
	log.InitLog(config)
	l := log.Entry()
	l.Println("xChain")
}
