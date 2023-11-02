package main

import (
	"github.com/develope/MonitorAgent/log"
	"github.com/develope/MonitorAgent/metrics"
	"github.com/develope/MonitorAgent/systeminfo"
)

func main() {
	go metrics.StartMetrics("node1", 9000)

	mon := systeminfo.SystemMonitor{}
	mon.Start()

	config := log.LogConfig{
		Kafka: &log.KafkaConfig{
			Brokers: []string{"127.0.0.1:9092"},
			//Brokers: []string{"47.104.157.94:9092"},
			Topic: "nodelog",
		},
		Path:  "logs",
		Level: "debug",
	}
	log.InitLog(config)
	l := log.Entry()
	l.Println("xChain")

	ch := make(chan interface{})
	<-ch
}
