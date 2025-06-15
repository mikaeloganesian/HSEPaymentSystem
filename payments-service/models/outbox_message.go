package models

import "github.com/google/uuid"

type OutboxMessage struct {
	ID        uuid.UUID
	EventType string
	Payload   []byte
}
