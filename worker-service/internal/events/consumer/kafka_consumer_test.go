package consumer

import (
	"errors"
	"fmt"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/blacklane/worker/testing/mocks"
)

func TestGetEventNameFromJSON(t *testing.T) {
	assert := assert.New(t)

	expectedEventName := "ImportantEvent"
	jsonString := generateEventJSON(expectedEventName)

	// act
	eventName, err := getEventNameFromJSON([]byte(jsonString))

	// assert
	assert.Nil(err)
	assert.Equal(eventName, expectedEventName)
}

func TestKafkaConsumer_ReceiveMessage_CallsHandler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := 0
	mockHandler, done := getMockHandler(&handlerCalled)
	mockKafkaWrapper := &mocks.MockKafkaReaderWrapper{}
	messageContent := generateEventJSON("anyEvent")
	message := kafka.Message{
		Topic: "test-topic",
		Value: []byte(messageContent),
	}
	consumer := KafkaConsumer{
		reader:    mockKafkaWrapper,
		listeners: []HandlerFunc{mockHandler},
	}

	mockKafkaWrapper.On("ReadMessage", mock.Anything).Return(message, nil).Once()

	// act
	consumer.receiveMessage()

	// wait for async go routine
	<-done

	// assert
	mockKafkaWrapper.AssertExpectations(t)
	assert.Equal(1, handlerCalled)
}

func TestKafkaConsumer_ReceiveMessage_HandlesError(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := 0
	mockHandler, _ := getMockHandler(&handlerCalled)
	mockKafkaWrapper := &mocks.MockKafkaReaderWrapper{}
	consumer := KafkaConsumer{
		reader:    mockKafkaWrapper,
		listeners: []HandlerFunc{mockHandler},
	}

	mockKafkaWrapper.On("ReadMessage", mock.Anything).Return(kafka.Message{}, errors.New("Error")).Once()

	// act
	consumer.receiveMessage()

	// assert
	mockKafkaWrapper.AssertExpectations(t)
	assert.Equal(0, handlerCalled)
}

func getMockHandler(callCounter *int) (HandlerFunc, <-chan bool) {
	done := make(chan bool, 1)
	return func(jsonBytes []byte, eventName string) {
		*callCounter++
		done <- true
	}, done
}

func generateEventJSON(eventName string) string {
	return fmt.Sprintf(`{
                "event": "%v",
                "payload": {
                        "someValue": "..."
                }
        }`, eventName)
}
