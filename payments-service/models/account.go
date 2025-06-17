package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        int       `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Balance   int64     `db:"balance" json:"balance"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
