package handlers

import (
	"net/http"
	"payments-service/models"

	"github.com/gin-gonic/gin"
)

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
