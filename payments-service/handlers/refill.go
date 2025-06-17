package handlers

import (
	"net/http"
	"payments-service/models"

	"github.com/gin-gonic/gin"
)

// RefillBalance godoc
// @Summary Пополнить баланс
// @Description Пополняет баланс счета
// @Tags accounts
// @Accept json
// @Produce json
// @Success 202 {object} map[string]string "Order accepted"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts/{id}/deposit [post]
func (h *Handler) RefillBalance(c *gin.Context) {
	var input models.RefillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.DB.Exec(`
        UPDATE accounts
        SET balance = balance + $1
        WHERE user_id = $2
    `, input.Amount, input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refill balance"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance refilled"})
}
