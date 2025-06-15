package main

import (
	"payments-service/db"
	"payments-service/handlers"
	"payments-service/outbox"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	outbox.StartOutboxPublisher(db.DB.DB)

	r := gin.Default()
	h := handlers.Handler{DB: db.DB}

	r.POST("/accounts", h.CreateAccount)
	r.POST("/accounts/:id/deposit", h.Deposit)
	r.GET("/accounts/:id/balance", h.GetBalance)

	r.Run(":8080")
}
