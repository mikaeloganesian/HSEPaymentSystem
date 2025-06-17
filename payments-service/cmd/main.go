package main

import (
	"payments-service/db"
	"payments-service/handlers"
	worker "payments-service/services"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func main() {
	db.InitDB()

	r := gin.Default()
	h := handlers.Handler{DB: db.DB}

	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "payment-status",
		Balancer: &kafka.LeastBytes{},
	}
	publisher := &worker.OutboxPublisher{
		DB:     db.DB.DB,
		Writer: writer,
	}

	go publisher.Start()

	go worker.NewWorker(db.DB).Start()

	r.POST("/accounts", h.CreateAccount)
	r.POST("/accounts/:id/deposit", h.RefillBalance)
	r.GET("/accounts/balance", h.GetBalance)

	r.Run(":8080")
}
