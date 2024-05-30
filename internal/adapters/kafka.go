package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
)

type KafkaAdapter struct {
	producer *kafka.Producer
}

func NewKafkaAdapter(kafkaHost string) (*KafkaAdapter, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:9092", kafkaHost),
	})
	if err != nil {
		return nil, err
	}

	return &KafkaAdapter{
		producer: producer,
	}, nil
}

func (adapter *KafkaAdapter) PublishChangelogEntry(entry entities.ChangelogEntry) error {
	topicName := "users-changelog"

	topicMessage, err := json.Marshal(entry)
	if err != nil {
		slog.Debug("unable to send changelog event", "err", err)
		return err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
		Value:          topicMessage,
		Key:            []byte(entry.UserID.String()),
	}

	deliveryChan := make(chan kafka.Event)
	err = adapter.producer.Produce(msg, nil)
	if err != nil {
		slog.Error("failed to produce message to DLQ", err)
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Failed to deliver message: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Produced message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)

	adapter.producer.Flush(1 * 1000)
	return nil
}
