package handlers

import (
	"net/http"
	"orders-service/models"

	"github.com/gin-gonic/gin"
)

// GetOrders godoc
// @Summary Получить список заказов
// @Description Возвращает список всех заказов
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /orders [get]
func (h *Handler) GetOrders(c *gin.Context) {
	var orders []models.Order

	err := h.DB.Select(&orders, "SELECT * FROM orders ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
