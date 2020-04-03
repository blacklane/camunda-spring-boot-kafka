package mocks

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

type MockKafkaWriterWrapper struct {
	mock.Mock
}

func (mkw *MockKafkaWriterWrapper) Close() error {
	return nil
}

func (mkw *MockKafkaWriterWrapper) WriteMessage(ctx context.Context, message kafka.Message) error {
	args := mkw.Called(ctx, message)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}

	return nil
}
