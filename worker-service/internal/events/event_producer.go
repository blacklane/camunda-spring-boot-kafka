package events

import "context"

type EventProducer interface {
	EmitRideAcceptedEvent(ctx context.Context, rideUUID string)
	Close() error
}


const (
	RideAccepted string = "RideAccepted"
)