package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blacklane/warsaw/logger"
	kafka "github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go/snappy" // snappy is required to decode compressed messages from kafka

	"github.com/blacklane/worker/internal/events"
)

type (
	HandlerFunc func(jsonBytes []byte, eventName string, writer events.EventProducer)

	KafkaConsumer struct {
		listeners []HandlerFunc
		reader    events.EventConsumer
		writer 	  events.EventProducer
	}
)

func newKafkaConsumer(reader events.EventConsumer, writer events.EventProducer) *KafkaConsumer {
	return &KafkaConsumer{
		listeners: []HandlerFunc{},
		reader:    reader,
		writer:	   writer,
	}
}

func StartKafkaConsumer(reader *kafka.Reader, listener HandlerFunc, writer events.EventProducer) {
	consumer := newKafkaConsumer(NewKafkaReaderWrapper(reader), writer)
	consumer.AddListener(listener)

	logger.Event("kafka_consumer_started").Msg(fmt.Sprintf("Starting kafka consumer on topic - %v\n", reader.Config().Topic))
	go consumer.Consume()
}

func (kafka *KafkaConsumer) AddListener(listener HandlerFunc) {
	kafka.listeners = append(kafka.listeners, listener)
}

func (kafka *KafkaConsumer) Consume() {
	defer kafka.reader.Close()

	for {
		kafka.receiveMessage()
	}
}

func (kafka *KafkaConsumer) receiveMessage() {
	m, err := kafka.reader.ReadMessage(context.Background())
	if err != nil {
		logger.Error("kafka_error", err).Msg(fmt.Sprintf("Consumer failed to read message: %+v", m))
		return
	}

	eventName, err := getEventNameFromJSON(m.Value)
	if err != nil {
		logger.Error("kafka_error", err).Msg(fmt.Sprintf("Couldn't read event name from message: %v", string(m.Value)))
		return
	}

	logger.Event("event_received").Msg(fmt.Sprintf("Received Event %v on topic %v", eventName, kafka.reader.Config().Topic))
	for _, handler := range kafka.listeners {
		go handler(m.Value, eventName, kafka.writer)
	}
}

func GetKafkaReader(brokers []string, topic string, consumerGroup string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  consumerGroup,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func UnmarshalJSONEvent(jsonBytes []byte, eventMap interface{}) {
	err := json.Unmarshal(jsonBytes, eventMap)
	if err != nil {
		logger.Error("kafka_error", err).Msg("Couldn't unmarshal kafka event.")
	}
}

func getEventNameFromJSON(jsonBytes []byte) (string, error) {
	var dat map[string]interface{}

	err := json.Unmarshal(jsonBytes, &dat)
	if err != nil {
		logger.Error("kafka_error", err).Msg("Event unmarsheling failed with error.")
		return "", err
	}

	if name, ok := dat["event"].(string); ok {
		return name, nil
	}

	return "", fmt.Errorf("event name not found")
}
