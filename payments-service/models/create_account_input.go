package models

import "github.com/google/uuid"

type CreateAccountInput struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
