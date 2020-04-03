package events

import "context"

type EventProducer interface {
	EmitEvent(ctx context.Context)
	Close() error
}
