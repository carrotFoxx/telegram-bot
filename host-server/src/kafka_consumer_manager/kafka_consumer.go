package kafka_consumer_manager

import (
	"Iot/src/common"
	"Iot/src/host_manager"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

/*
Модуль отвечает за получение сообщений для отправления комманд хостам.
Формат входященго сообщения:
{
"Id": int64, - id хоста
"task": int32, - номер команды
"arg": string - аргумент команды
}
*/

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaURL},
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       1,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})
}

func Kafka_consumer_service() {

	kafkaURL := os.Getenv("KAFKA_DSN")
	topic := os.Getenv("INCOMING_TOPIC")
	groupID := os.Getenv("GROUP_ID")
	reader := getKafkaReader(kafkaURL, topic, groupID)

	defer reader.Close()

	log.Println("start consuming ... !!")
	msg := common.Command_msg{}

	for {
		kafka_msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(kafka_msg.Value, &msg)
		if err != nil {
			log.Println(err)
		} else {
			host_manager.Send_message_from_outside(msg)
		}
		fmt.Println(msg)
		fmt.Println(string(kafka_msg.Value))
	}
}
