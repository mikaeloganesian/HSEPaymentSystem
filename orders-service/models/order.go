package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderInput struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	Amount      int64     `json:"amount" binding:"required"`
	Description string    `json:"description"`
	Status      string    `json:"status" binding:"required"`
}

type Order struct {
	ID          int       `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Amount      int64     `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
