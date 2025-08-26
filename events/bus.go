package events

import (
	"encoding/json"
	"errors"
	"github.com/conv-project/go-shared/kafka"
)

var (
	ErrEventMarshal = errors.New("validation error")
)

type EventBus struct {
	producer *kafka.Producer
	topic    string
	group    string
}

func NewEventBus(producer *kafka.Producer, topic string, group string) *EventBus {
	return &EventBus{
		producer: producer,
		topic:    topic,
		group:    group,
	}
}

func (e *EventBus) Publish(actorId, eventType string, payload interface{}) error {
	event := NewEvent(
		actorId, eventType, e.group, payload,
	)
	if event == nil {
		return ErrEventMarshal
	}

	marshalled, err := json.Marshal(event)
	if err != nil {
		return ErrEventMarshal
	}

	return e.producer.SendMessage(e.topic, actorId, marshalled)
}
