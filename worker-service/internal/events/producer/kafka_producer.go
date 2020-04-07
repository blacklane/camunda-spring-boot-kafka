package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blacklane/warsaw/logger"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"

	"github.com/blacklane/worker/internal/events"
)

type KafkaProducer struct {
	writer KafkaWriter
}

func NewKafkaWriter(brokers []string, topic string) events.EventProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            topic,
		CompressionCodec: snappy.NewCompressionCodec(),
	})

	return &KafkaProducer{
		writer: NewKafkaWriterWrapper(writer),
	}
}

func (producer *KafkaProducer) EmitRideAcceptedEvent(ctx context.Context, rideUUID string) {
	event := &RideAccepted{
		Event: "RideAccepted",
		Payload: rideUUID,
	}
	messageKey := fmt.Sprintf("ride-accepted-%s", rideUUID)

	jsonValue, err := json.Marshal(event)
	if err != nil {
		logger.Error("kafka_error", err).Msg(fmt.Sprintf("Failed to marshal event %v", event))
		return
	}

	producer.emitMessage(ctx, messageKey, jsonValue)
}

func (producer *KafkaProducer) emitMessage(ctx context.Context, messageKey string, jsonMessage []byte) {
	message := kafka.Message{
		Key:   []byte(messageKey),
		Value: jsonMessage,
	}

	err := producer.writer.WriteMessage(ctx, message)
	if err != nil {
		logger.Error("kafka_error", err).Msg(fmt.Sprintf("Failed to emit event %s", messageKey))
	}

	logger.Event("event_emitted").Msg(fmt.Sprintf("Emitted event %s", messageKey))
}

func (producer *KafkaProducer) Close() error {
	return producer.writer.Close()
}
