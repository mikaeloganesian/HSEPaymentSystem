// services/publisher.go
package services

import (
	"context"
	"encoding/json"
	"time"

	"orders-service/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Publisher struct {
	DB *sqlx.DB
}

func (p *Publisher) CreateOrder(ctx context.Context, order models.Order) error {
	tx, err := p.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Вставляем заказ
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	_, err = tx.NamedExec(`INSERT INTO orders (id, user_id, amount, description, status, created_at, updated_at)
        VALUES (:id, :user_id, :amount, :description, :status, :created_at, :updated_at)`, &order)
	if err != nil {
		return err
	}

	// Готовим Payload
	payload, err := json.Marshal(order)
	if err != nil {
		return err
	}

	outbox := models.Outbox{
		ID:        uuid.New().String(),
		EventType: "order_created",
		Payload:   payload,
		CreatedAt: time.Now(),
		Sent:      false,
	}

	_, err = tx.NamedExec(`INSERT INTO outbox (id, event_type, payload, created_at, sent)
        VALUES (:id, :event_type, :payload, :created_at, :sent)`, &outbox)
	if err != nil {
		return err
	}

	return tx.Commit()
}
