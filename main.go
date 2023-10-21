package main

import (
	"github.com/develope/MonitorAgent/log"
	"github.com/develope/MonitorAgent/metrics"
	"github.com/develope/MonitorAgent/systeminfo"
)

func main() {
	go metrics.StartMetrics(9000)

	mon := systeminfo.SystemMonitor{}
	mon.Start()

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

	ch := make(chan interface{})
	<-ch
}
