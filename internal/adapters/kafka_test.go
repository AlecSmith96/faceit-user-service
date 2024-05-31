package adapters_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AlecSmith96/faceit-user-service/internal/adapters"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	mock_adapters "github.com/AlecSmith96/faceit-user-service/mocks/adapters"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"github.com/segmentio/kafka-go"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
	"time"
)

var (
	ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func TestNewKafkaAdapter(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	mockDialer := mock_adapters.NewMockDialer(ctrl)
	mockKafkaConnection := mock_adapters.NewMockKafkaConnection(ctrl)

	mockDialer.EXPECT().
		DialLeader(gomock.AssignableToTypeOf(ctxType), "tcp", "localhost", "users-changelog", 0).
		Return(mockKafkaConnection, nil)

	adapter, err := adapters.NewKafkaAdapter("localhost", mockDialer)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(adapter).To(BeAssignableToTypeOf(&adapters.KafkaAdapter{}))
}

func TestNewKafkaAdapter_ReturnsErr(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	mockDialer := mock_adapters.NewMockDialer(ctrl)

	mockDialer.EXPECT().
		DialLeader(gomock.AssignableToTypeOf(ctxType), "tcp", "localhost", "users-changelog", 0).
		Return(nil, errors.New("an error occurred"))

	adapter, err := adapters.NewKafkaAdapter("localhost", mockDialer)
	g.Expect(err).To(MatchError("an error occurred"))
	g.Expect(adapter).To(BeNil())
}

func TestKafkaAdapter_PublishChangelogEntry(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	mockDialer := mock_adapters.NewMockDialer(ctrl)
	mockKafkaConnection := mock_adapters.NewMockKafkaConnection(ctrl)

	entry := entities.ChangelogEntry{
		UserID:     uuid.New(),
		CreatedAt:  time.Now(),
		ChangeType: "POST",
	}
	entryJSON, err := json.Marshal(entry)
	g.Expect(err).ToNot(HaveOccurred())

	mockDialer.EXPECT().
		DialLeader(gomock.AssignableToTypeOf(ctxType), "tcp", "localhost", "users-changelog", 0).
		Return(mockKafkaConnection, nil)
	mockKafkaConnection.EXPECT().SetWriteDeadline(gomock.AssignableToTypeOf(time.Time{})).Return(nil)
	mockKafkaConnection.EXPECT().WriteMessages(kafka.Message{Value: entryJSON}).Return(0, nil)

	adapter, err := adapters.NewKafkaAdapter("localhost", mockDialer)

	err = adapter.PublishChangelogEntry(entry)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestKafkaAdapter_PublishChangelogEntry_SetWriteDeadlineErr(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	mockDialer := mock_adapters.NewMockDialer(ctrl)
	mockKafkaConnection := mock_adapters.NewMockKafkaConnection(ctrl)

	entry := entities.ChangelogEntry{
		UserID:     uuid.New(),
		CreatedAt:  time.Now(),
		ChangeType: "POST",
	}

	mockDialer.EXPECT().
		DialLeader(gomock.AssignableToTypeOf(ctxType), "tcp", "localhost", "users-changelog", 0).
		Return(mockKafkaConnection, nil)
	mockKafkaConnection.EXPECT().SetWriteDeadline(gomock.AssignableToTypeOf(time.Time{})).Return(errors.New("an error occurred"))

	adapter, err := adapters.NewKafkaAdapter("localhost", mockDialer)

	err = adapter.PublishChangelogEntry(entry)
	g.Expect(err).To(MatchError("an error occurred"))
}

func TestKafkaAdapter_WriteMessagesErr(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	mockDialer := mock_adapters.NewMockDialer(ctrl)
	mockKafkaConnection := mock_adapters.NewMockKafkaConnection(ctrl)

	entry := entities.ChangelogEntry{
		UserID:     uuid.New(),
		CreatedAt:  time.Now(),
		ChangeType: "POST",
	}
	entryJSON, err := json.Marshal(entry)
	g.Expect(err).ToNot(HaveOccurred())

	mockDialer.EXPECT().
		DialLeader(gomock.AssignableToTypeOf(ctxType), "tcp", "localhost", "users-changelog", 0).
		Return(mockKafkaConnection, nil)
	mockKafkaConnection.EXPECT().SetWriteDeadline(gomock.AssignableToTypeOf(time.Time{})).Return(nil)
	mockKafkaConnection.EXPECT().WriteMessages(kafka.Message{Value: entryJSON}).Return(0, errors.New("an error occurred"))

	adapter, err := adapters.NewKafkaAdapter("localhost", mockDialer)

	err = adapter.PublishChangelogEntry(entry)
	g.Expect(err).To(MatchError("an error occurred"))
}

func TestKafkaAdapter_CloseConn(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	mockDialer := mock_adapters.NewMockDialer(ctrl)
	mockKafkaConnection := mock_adapters.NewMockKafkaConnection(ctrl)

	mockDialer.EXPECT().
		DialLeader(gomock.AssignableToTypeOf(ctxType), "tcp", "localhost", "users-changelog", 0).
		Return(mockKafkaConnection, nil)
	mockKafkaConnection.EXPECT().Close().Return(nil)

	adapter, err := adapters.NewKafkaAdapter("localhost", mockDialer)
	g.Expect(err).ToNot(HaveOccurred())

	err = adapter.CloseConn()
	g.Expect(err).ToNot(HaveOccurred())
}

func TestKafkaAdapter_CloseConnErr(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	mockDialer := mock_adapters.NewMockDialer(ctrl)
	mockKafkaConnection := mock_adapters.NewMockKafkaConnection(ctrl)

	mockDialer.EXPECT().
		DialLeader(gomock.AssignableToTypeOf(ctxType), "tcp", "localhost", "users-changelog", 0).
		Return(mockKafkaConnection, nil)
	mockKafkaConnection.EXPECT().Close().Return(errors.New("an error occurred"))

	adapter, err := adapters.NewKafkaAdapter("localhost", mockDialer)
	g.Expect(err).ToNot(HaveOccurred())

	err = adapter.CloseConn()
	g.Expect(err).To(MatchError("an error occurred"))
}
