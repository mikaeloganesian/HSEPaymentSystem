// consumer/order_consumer.go
package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type Order struct {
	UserID      uuid.UUID `json:"user_id"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
}

func StartOrderConsumer(db *sqlx.DB) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "orders",
		GroupID:  "order-consumer-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	log.Println("Order consumer started...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var order Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Failed to parse order: %v", err)
			continue
		}

		now := time.Now()

		query := `
			INSERT INTO orders (user_id, amount, description, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

		_, err = db.Exec(query, order.UserID, order.Amount, order.Description, order.Status, now, now)
		if err != nil {
			log.Printf("Failed to insert order: %v", err)
		} else {
			log.Printf("Order inserted: user_id=%s amount=%d", order.UserID, order.Amount)
		}
	}
}
