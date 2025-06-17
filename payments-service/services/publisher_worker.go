package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type OutboxPublisher struct {
	DB     *sql.DB
	Writer *kafka.Writer
}

type PaymentStatusEvent struct {
	OrderID string `json:"id"`
	Status  string `json:"status"`
}

func (p *OutboxPublisher) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rows, err := p.DB.Query(`
			SELECT id, event_type, payload 
			FROM outbox 
			WHERE sent = false 
			AND event_type IN ('PaymentSuccess', 'PaymentFailed')
			ORDER BY created_at 
			LIMIT 10
		`)
		if err != nil {
			log.Println("Failed to fetch outbox events:", err)
			continue
		}

		for rows.Next() {
			var id string
			var eventType string
			var payloadBytes []byte

			if err := rows.Scan(&id, &eventType, &payloadBytes); err != nil {
				log.Println("Failed to scan row:", err)
				continue
			}

			var event PaymentStatusEvent
			if err := json.Unmarshal(payloadBytes, &event); err != nil {
				log.Printf("Failed to unmarshal payload: %v", err)
				continue
			}

			// Устанавливаем статус на основе event_type
			switch eventType {
			case "PaymentSuccess":
				event.Status = "success"
			case "PaymentFailed":
				event.Status = "failed"
			default:
				log.Printf("Unknown event type: %s", eventType)
				continue
			}

			// Кодируем новый payload
			encoded, err := json.Marshal(event)
			if err != nil {
				log.Printf("Failed to re-encode event: %v", err)
				continue
			}

			log.Printf("Publishing payment status event to Kafka: id=%s status=%s", event.OrderID, event.Status)

			err = p.Writer.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte(event.OrderID),
					Value: encoded,
				})
			if err != nil {
				log.Println("Failed to publish message to Kafka:", err)
				continue
			}

			_, err = p.DB.Exec(`UPDATE outbox SET sent = true WHERE id = $1`, id)
			if err != nil {
				log.Println("Failed to update outbox status:", err)
			}
		}
		rows.Close()
	}
}
