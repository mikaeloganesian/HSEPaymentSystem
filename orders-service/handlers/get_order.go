package handlers

import (
	"net/http"
	"orders-service/models"

	"github.com/gin-gonic/gin"
)

// GetOrderStatus godoc
// @Summary Получить статус заказа
// @Description Возвращает статус заказа по ID
// @Tags orders
// @Produce json
// @Param id path int true "ID заказа"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /orders/{id}/status [get]
func (h *Handler) GetOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	err := h.DB.Get(&order, "SELECT id, status FROM orders WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": order.ID, "status": order.Status})
}
