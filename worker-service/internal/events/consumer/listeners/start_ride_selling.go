package listeners

import (
	"github.com/blacklane/warsaw/logger"
	"github.com/blacklane/worker/internal/events"
	"github.com/blacklane/worker/internal/events/consumer"
	"context"
)

func StartRideSellingListener(jsonBytes []byte, eventName string, writer events.EventProducer) {
	switch eventName {
	case events.StartRideSelling:
		var event consumer.StartRideSellingEvent
		consumer.UnmarshalJSONEvent(jsonBytes, &event)
		handleStartRideSelling(event, writer)
	}
}

func handleStartRideSelling(event consumer.StartRideSellingEvent, writer events.EventProducer) {
	logger.Event(event.Event).Msg("Starting ride selling")

	// defer 20 seconds

	// send accepted message
	writer.EmitRideAcceptedEvent(context.Background(), event.Payload)
}