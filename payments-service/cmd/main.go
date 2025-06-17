// @title Payment Service API
// @version 1.0
// @description API для управления счетами
// @host localhost:8080
// @BasePath /
package main

import (
	swaggerFiles "github.com/swaggo/files"

	"payments-service/db"
	"payments-service/handlers"
	worker "payments-service/services"

	_ "payments-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
