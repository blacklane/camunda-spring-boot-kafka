package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/blacklane/worker/internal/events"
)

type KafkaReaderWrapper struct {
	reader *kafka.Reader
}

func (mr *KafkaReaderWrapper) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return mr.reader.ReadMessage(ctx)
}

func (mr *KafkaReaderWrapper) Close() error {
	return mr.reader.Close()
}

func (mr *KafkaReaderWrapper) Config() kafka.ReaderConfig {
	return mr.reader.Config()
}

func NewKafkaReaderWrapper(reader *kafka.Reader) events.EventConsumer {
	return &KafkaReaderWrapper{reader: reader}
}
