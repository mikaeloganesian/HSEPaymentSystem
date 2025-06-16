package models

import "github.com/google/uuid"

type OutboxEvent struct {
	ID        uuid.UUID
	EventType string
	Payload   interface{}
}
