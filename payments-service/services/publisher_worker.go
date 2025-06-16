package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type OutboxPublisher struct {
	DB     *sql.DB
	Writer *kafka.Writer
}

func (p *OutboxPublisher) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rows, err := p.DB.Query(`
            SELECT id, event_type, payload 
            FROM outbox 
            WHERE sent = false
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
				continue
			}

			// Публикуем в Kafka
			err = p.Writer.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte(id),
					Value: payloadBytes,
				})
			if err != nil {
				log.Println("Failed to publish message to Kafka:", err)
				continue
			}

			// Обновляем запись
			_, err = p.DB.Exec(`UPDATE outbox SET sent = true WHERE id = $1`, id)
			if err != nil {
				log.Println("Failed to update outbox status:", err)
			}
		}
		rows.Close()
	}
}
