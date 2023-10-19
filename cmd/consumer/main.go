package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

var (
	broker     = flag.String("broker", "localhost:9092", "The Kafka broker address")
	topic      = flag.String("topic", "topic_1", "Topic to consumer")
	partition  = flag.Int("p", 0, "partition to consumer")
	offsetType = flag.Int("offsetType", 0, "Offset Type (OffsetNewest | OffsetOldest)")
)

func main() {
	flag.Parse()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := []string{*broker}
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Connected to Kafka cluster %v\n", brokers)
	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(*topic, int32(*partition), sarama.OffsetOldest)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Consuming topic %v partition %v\n", *topic, *partition)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				log.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed", "messages")
}
