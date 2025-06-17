package service

import (
	"context"
	"encoding/json"
	"log"

	"payments-service/models"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type Worker struct {
	DB     *sqlx.DB
	Reader *kafka.Reader
}

func NewWorker(db *sqlx.DB) *Worker {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "orders",
		GroupID:  "order-service",
		MinBytes: 1e3,
		MaxBytes: 1e6,
	})

	return &Worker{
		DB:     db,
		Reader: r,
	}
}
func (w *Worker) Start() {
	log.Println("[worker] Starting consumer loop")
	for {
		m, err := w.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[worker] Error reading message: %v", err)
			continue
		}
		log.Printf("[worker] Received message: topic=%s partition=%d offset=%d key=%s value=%s",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		var event models.OrderEvent
		err = json.Unmarshal(m.Value, &event)
		if err != nil {
			log.Printf("[worker] Failed to unmarshal message: %v", err)
			continue
		}

		err = w.handleOrder(event)
		if err != nil {
			log.Printf("[worker] Error handling order: %v", err)
		}
	}
}
