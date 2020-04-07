package producer

type (
	RideAccepted struct {
		Event     string `json:"event"`
		Payload   string `json:"payload"`
	}
)
