package models

import "github.com/google/uuid"

type RefillInput struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Amount int64     `json:"amount" binding:"required,gt=0"`
}
