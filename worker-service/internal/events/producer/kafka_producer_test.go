package producer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/blacklane/worker/testing/mocks"
)

func TestKafkaProducer_EmitChatRoomCreated(t *testing.T) {
	mockKafkaWriterWrapper := &mocks.MockKafkaWriterWrapper{}
	producer := &KafkaProducer{
		writer: mockKafkaWriterWrapper,
	}

	mockKafkaWriterWrapper.On("WriteMessage", mock.Anything, mock.Anything).Return(nil)

	// act
	producer.EmitEvent(context.Background())

	// assert
	mockKafkaWriterWrapper.AssertExpectations(t)
}
