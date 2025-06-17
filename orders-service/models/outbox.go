package models

import "time"

type Outbox struct {
	ID        string    `db:"id"`
	EventType string    `db:"event_type"`
	Payload   []byte    `db:"payload"`
	CreatedAt time.Time `db:"created_at"`
	Sent      bool      `db:"sent"`
}
