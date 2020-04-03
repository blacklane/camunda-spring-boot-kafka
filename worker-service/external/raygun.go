package external

import "github.com/MindscapeHQ/raygun4go"

type ErrorReporter interface {
	SendError(error error)
}

type Raygun struct {
	client *raygun4go.Client
}

type noop struct{}

// NewRaygun returns a ErrorReporter reporting to Raygun
// If client is nil, it'll be a no-op
func NewRaygun(client *raygun4go.Client) ErrorReporter {
	if client == nil {
		return noop{}
	}

	return Raygun{client: client}
}

func (n noop) SendError(err error) {}

func (r Raygun) SendError(err error) {
	r.client.SendError(err)
}
