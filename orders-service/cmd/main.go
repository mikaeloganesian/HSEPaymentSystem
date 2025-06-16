// @title Orders Service API
// @version 1.0
// @description API для управления заказами
// @host localhost:8081
// @BasePath /
package main

import (
	"context"
	"fmt"
	"orders-service/db"
	"orders-service/handlers"
	"orders-service/services"

	_ "orders-service/docs"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db.InitDB()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	})

	publisher := &services.Publisher{DB: db.DB}
	worker := &services.OutboxWorker{DB: db.DB, Writer: writer}
	handler := &handlers.Handler{Publisher: publisher, DB: db.DB}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler.StartPaymentStatusConsumer(ctx)

	go worker.Run(context.Background())

	r := gin.Default()
	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders", handler.GetOrders)
	r.GET("/order/:id", handler.GetOrderStatus)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8081")

	fmt.Println("swagger documentation address: http://localhost:8081/swagger/index.html")
}
