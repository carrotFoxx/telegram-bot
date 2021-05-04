package kafka_producer_manager

import (
	"context"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"log"
	"os"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

var writer = newKafkaWriter(os.Getenv("KAFKA_DSN"), os.Getenv("OUTGOING_TOPIC"))

func Send_kafka_msg(bts []byte) {
	msg := kafka.Message{
		Value: bts,
	}
	log.Println("incomming_info_message: %s", string(bts))
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
	}
}
