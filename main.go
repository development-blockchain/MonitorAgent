package main

import (
	"github.com/develope/MonitorAgent/log"
	"github.com/develope/MonitorAgent/metrics"
	"github.com/develope/MonitorAgent/systeminfo"
	"math/big"
	"math/rand"
	"time"
)

func demoNode() {
	l := log.Entry()
	height := big.NewInt(0)
	one := big.NewInt(1)
	t := time.NewTicker(time.Second * 3)
	methods := []string{"getBalance", "getBlock", "getAccount"}
	defer t.Stop()
	for {
		select {
		case <-t.C:
			height = height.Add(height, one)
			metrics.UpdateChainHeight(height)
			metrics.UpdatePeerCount(rand.Intn(10))
			idx := int(height.Int64()) % len(methods)
			metrics.UpdateRestHandler(methods[idx])
			l.Info("demo node info")
		}
	}
}

func main() {
	go metrics.StartMetrics("node1", 9000)

	mon := systeminfo.SystemMonitor{}
	mon.Start()

	config := log.LogConfig{
		//Kafka: &log.KafkaConfig{
		//	Brokers: []string{"127.0.0.1:9092"},
		//	//Brokers: []string{"47.104.157.94:9092"},
		//	Topic: "nodelog",
		//},
		Path:  "logs",
		Level: "debug",
	}
	log.InitLog(config)
	l := log.Entry()
	l.Println("xChain")
	go demoNode()

	ch := make(chan interface{})
	<-ch
}
