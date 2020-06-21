package eventpublisher

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type EventPublisher struct{
	nc *nats.Conn
}

func NewEventPublisher(nc *nats.Conn) *EventPublisher {
	return &EventPublisher{nc: nc}
}

type AuthSuccessEventMessage struct {
	ID string `json:"id"`
}

func (e *EventPublisher) PublishNewAuthSuccessEvent(ctx context.Context, userID string) {
	message := AuthSuccessEventMessage{ID: userID}
	bs, _ := json.Marshal(message)
	e.nc.Publish("auth.creation.success", bs)
}

func (e *EventPublisher) PublishNewAuthFailedEvent(ctx context.Context, userID string) {
	message := AuthSuccessEventMessage{ID: userID}
	bs, _ := json.Marshal(message)
	e.nc.Publish("auth.creation.failed", bs)
}