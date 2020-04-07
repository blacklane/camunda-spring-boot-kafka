package consumer

type (
	StartRideSellingEvent struct {
		Event   string `json:"event"`
		Payload string `json:"payload"`
	}
)
