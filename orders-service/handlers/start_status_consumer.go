package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func (h *Handler) StartPaymentStatusConsumer(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "payment-status",
		GroupID:  "order-service-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(ctx)
			fmt.Printf("fint topic message!!!")
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Printf("error reading kafka message: %v", err)
				continue
			}

			var event PaymentStatusEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("failed to unmarshal payment event: %v", err)
				continue
			}

			// Обновляем статус заказа
			if err := h.updateOrderStatus(event.OrderID, event.Status); err != nil {
				log.Printf("failed to update order status: %v", err)
			}
		}
	}()
}

type PaymentStatusEvent struct {
	OrderID string `json:"ID"`
	Status  string `json:"Status"`
}

func (h *Handler) updateOrderStatus(orderID string, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := h.DB.Exec(query, status, orderID)
	return err
}
