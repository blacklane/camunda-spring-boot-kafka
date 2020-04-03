package producer

import "time"

type (
	SomeEvent struct {
		Event     string           `json:"event"`
		CreatedAt time.Time        `json:"created_at"`
		Payload   SomeEventPayload `json:"payload"`
	}

	SomeEventPayload struct {
		SomeValue string `json:"some_value"`
	}
)
