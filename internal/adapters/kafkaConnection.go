package adapters

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

// KafkaConnection is an interface used for mocking kafka calls in tests
//
//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/adapters/kafkaConnection.go  . "KafkaConnection"
type KafkaConnection interface {
	SetWriteDeadline(t time.Time) error
	WriteMessages(msgs ...kafka.Message) (int, error)
	Close() error
}

// KafkaConnectionWrapper is used to wrap the underlying kafka library. It's functions explicitly cant be tested as they
// are replaced by the mocking functionality.
type KafkaConnectionWrapper struct {
	conn *kafka.Conn
}

func (kc *KafkaConnectionWrapper) SetWriteDeadline(t time.Time) error {
	return kc.conn.SetWriteDeadline(t)
}

func (kc *KafkaConnectionWrapper) WriteMessages(msgs ...kafka.Message) (int, error) {
	return kc.conn.WriteMessages(msgs...)
}

func (kc *KafkaConnectionWrapper) Close() error {
	return kc.conn.Close()
}

// Dialer is an interface used for mocking the connection to kafka in tests
//
//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/adapters/dialer.go  . "Dialer"
type Dialer interface {
	DialLeader(ctx context.Context, network string, kafkaHost string, topic string, partition int) (KafkaConnection, error)
}

// KafkaDialer is a struct that implements the Dialer interface. This allows for the mocking of the call to connect to
// kafka in tests.
type KafkaDialer struct{}

func (dialer *KafkaDialer) DialLeader(ctx context.Context, network string, kafkaHost string, topic string, partition int) (KafkaConnection, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", fmt.Sprintf("%s:9092", kafkaHost), topicName, 0)
	return &KafkaConnectionWrapper{
		conn: conn,
	}, err
}
