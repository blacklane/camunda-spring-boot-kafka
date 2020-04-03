package producer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type (
	KafkaWriter interface {
		Close() error
		WriteMessage(ctx context.Context, message kafka.Message) error
	}

	KafkaWriterWrapper struct {
		writer *kafka.Writer
	}
)

func NewKafkaWriterWrapper(writer *kafka.Writer) KafkaWriter {
	return &KafkaWriterWrapper{
		writer: writer,
	}
}

func (kww *KafkaWriterWrapper) Close() error {
	return kww.writer.Close()
}

func (kww *KafkaWriterWrapper) WriteMessage(ctx context.Context, message kafka.Message) error {

	return kww.writer.WriteMessages(ctx, message)
}
