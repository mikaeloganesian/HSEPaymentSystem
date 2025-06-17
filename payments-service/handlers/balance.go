package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetBalance godoc
// @Summary Создать заказ
// @Description Создает новый заказ и ставит в очередь
// @Tags accounts
// @Accept json
// @Produce json
// @Param order body models.Account true "Данные заказа"
// @Success 202 {object} map[string]string "Order accepted"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /account/{id}/balance [get]
func (h *Handler) GetBalance(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	var balance int64
	err := h.DB.Get(&balance, "SELECT balance FROM accounts WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
