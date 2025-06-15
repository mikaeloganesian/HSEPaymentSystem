package main

import (
	"payments-service/db"
	"payments-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()
	h := handlers.Handler{DB: db.DB}

	r.POST("/accounts", h.CreateAccount)
	r.POST("/accounts/:id/deposit", h.Deposit)
	r.GET("/accounts/:id/balance", h.GetBalance)

	r.Run(":8080")
}
