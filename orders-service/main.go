package main

import (
	"orders-service/consumer"
	"orders-service/db"
	"orders-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	go consumer.StartOrderConsumer(db.DB)

	r := gin.Default()
	h := handlers.Handler{DB: db.DB}

	r.POST("/orders", h.CreateOrder)

	r.Run(":8081")
}
