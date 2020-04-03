package mocks

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

type MockKafkaReaderWrapper struct {
	mock.Mock
}

func (mkw *MockKafkaReaderWrapper) ReadMessage(ctx context.Context) (kafka.Message, error) {
	args := mkw.Called(ctx)
	if args.Get(1) == nil {
		return args.Get(0).(kafka.Message), nil
	}
	return args.Get(0).(kafka.Message), args.Get(1).(error)
}

func (mkw *MockKafkaReaderWrapper) Close() error {
	return nil
}

func (mkw *MockKafkaReaderWrapper) Config() kafka.ReaderConfig {
	return kafka.ReaderConfig{}
}
