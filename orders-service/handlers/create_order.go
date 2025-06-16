// handlers/create_order.go
package handlers

import (
	"log"
	"net/http"
	"orders-service/models"
	"orders-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	Publisher *services.Publisher
	DB        *sqlx.DB
}

// CreateOrder godoc
// @Summary Создать заказ
// @Description Создает новый заказ и ставит в очередь
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.Outbox true "Данные заказа"
// @Success 202 {object} map[string]string "Order accepted"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	var input struct {
		UserID      uuid.UUID `json:"user_id"`
		Amount      int64     `json:"amount"`
		Description string    `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		ID:          uuid.New(),
		UserID:      input.UserID,
		Amount:      input.Amount,
		Description: input.Description,
		Status:      "created",
	}

	if err := h.Publisher.CreateOrder(c.Request.Context(), order); err != nil {
		log.Printf("error creating order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Order accepted"})
}
