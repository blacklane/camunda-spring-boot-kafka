package events

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type EventConsumer interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
	Config() kafka.ReaderConfig
}

const (
	StartRideSelling string = "StartRideSelling"
)