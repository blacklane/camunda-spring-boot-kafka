package consumer

type (
	SomeEvent struct {
		Event   string           `json:"event"`
		Payload SomeEventPayload `json:"payload"`
	}

	SomeEventPayload struct {
		SomeValue string `json:"some_value"`
	}
)
