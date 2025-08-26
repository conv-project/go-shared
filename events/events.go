package events

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

const (
	DefaultGroup = "default"
)

type Event struct {
	Id      string `json:"id"`
	ActorId string `json:"actor_id"`
	Type    string `json:"type"`
	Group   string `json:"group"`
	Ts      int64  `json:"ts"`
	Payload []byte `json:"payload"`
}

func NewEvent(actorId, eventType, group string, payload interface{}) *Event {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil
	}

	return &Event{
		Id:      uuid.New().String(),
		ActorId: actorId,
		Type:    eventType,
		Group:   group,
		Ts:      time.Now().UTC().Unix(),
		Payload: p,
	}
}
