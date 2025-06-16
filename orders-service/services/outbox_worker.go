// services/outbox_worker.go
package services

import (
	"context"
	"orders-service/models"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type OutboxWorker struct {
	DB     *sqlx.DB
	Writer *kafka.Writer
}

func (w *OutboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processOutbox(ctx)
		}
	}
}

func (w *OutboxWorker) processOutbox(ctx context.Context) {
	var messages []models.Outbox
	err := w.DB.SelectContext(ctx, &messages, `
        SELECT * FROM outbox WHERE sent = false AND event_type = 'order_created' LIMIT 10
    `)
	if err != nil || len(messages) == 0 {
		return
	}

	for _, msg := range messages {
		err := w.Writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(msg.ID),
			Value: msg.Payload,
		})
		if err != nil {
			continue // логируй ошибку
		}

		_, err = w.DB.ExecContext(ctx, `
            UPDATE outbox SET sent = true WHERE id = $1
        `, msg.ID)
		if err != nil {
			continue // логируй ошибку
		}
	}
}
