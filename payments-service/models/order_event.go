package models

import "time"

type OrderEvent struct {
	ID          string    `json:"ID"`
	Amount      int64     `json:"Amount"`
	Status      string    `json:"Status"`
	UserID      string    `json:"UserID"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	Description string    `json:"Description"`
}
