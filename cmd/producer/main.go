package main

import (
	"flag"
	"github.com/IBM/sarama"
	"log"
)

var (
	broker    = flag.String("broker", "localhost:9092", "The Kafka broker address")
	topic     = flag.String("topic", "topic1", "Topic to consumer")
	partition = flag.Int("p", 0, "partition to consumer")
)

func main() {
	flag.Parse()
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	brokerList := []string{*broker}
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()
	msg := &sarama.ProducerMessage{
		Topic:     *topic,
		Value:     sarama.StringEncoder("Something Cool"),
		Partition: int32(*partition),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", *topic, partition, offset)
}
