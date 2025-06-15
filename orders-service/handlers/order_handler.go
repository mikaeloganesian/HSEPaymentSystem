package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"orders-service/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	DB *sqlx.DB
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var input models.CreateOrderInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// сериализуем заказ
	data, err := json.Marshal(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode message"})
		return
	}

	// создаём writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(input.UserID.String()),
			Value: data,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message to Kafka"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Order accepted"})
}
