package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/segmentio/kafka-go"
	"log"
	"log/slog"
	"time"
)

const (
	topicName = "users-changelog"
)

type KafkaAdapter struct {
	conn *kafka.Conn
}

func NewKafkaAdapter(kafkaHost string) (*KafkaAdapter, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", fmt.Sprintf("%s:9092", kafkaHost), topicName, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	return &KafkaAdapter{
		conn: conn,
	}, nil
}

func (adapter *KafkaAdapter) PublishChangelogEntry(entry entities.ChangelogEntry) error {
	err := adapter.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		slog.Debug("unable to set write deadline", "err", err)
	}

	entryJSON, err := json.Marshal(entry)
	if err != nil {
		slog.Debug("unable to convert entry to json", "err", err)
		return err
	}

	_, err = adapter.conn.WriteMessages(
		kafka.Message{Value: entryJSON},
	)
	if err != nil {
		slog.Debug("failed to write message", "err", err)
		return err
	}

	return nil
}

func (adapter *KafkaAdapter) CloseConn() error {
	if err := adapter.conn.Close(); err != nil {
		slog.Error("failed to close writer", "err", err)
		return err
	}

	return nil
}
