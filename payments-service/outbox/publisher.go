package outbox

import (
	"database/sql"
	"log"
	"payments-service/models"
	"time"
)

func StartOutboxPublisher(db *sql.DB) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			rows, err := db.Query(`
				SELECT id, event_type, payload 
				FROM outbox 
				WHERE sent = FALSE
			`)
			if err != nil {
				log.Println("[Outbox] Error reading outbox:", err)
				continue
			}

			var messages []models.OutboxMessage
			for rows.Next() {
				var msg models.OutboxMessage
				if err := rows.Scan(&msg.ID, &msg.EventType, &msg.Payload); err != nil {
					log.Println("[Outbox] Scan error:", err)
					continue
				}
				messages = append(messages, msg)
			}
			rows.Close()

			for _, msg := range messages {
				log.Printf("[Outbox] Sending event: %s\nPayload: %s\n", msg.EventType, string(msg.Payload))

				// Отмечаем как отправленное
				_, err := db.Exec(`UPDATE outbox SET sent = TRUE WHERE id = $1`, msg.ID)
				if err != nil {
					log.Println("[Outbox] Mark sent failed:", err)
				}
			}
		}
	}()
}
