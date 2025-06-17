package handlers

import (
	"net/http"
	"payments-service/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetBalamce godoc
// @Summary Создать счет
// @Description Создает новый счет
// @Tags accounts
// @Accept json
// @Produce json
// @Success 202 {object} map[string]string "Order accepted"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var input models.CreateAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// check if account already exists
	var count int
	err := h.DB.Get(&count, "SELECT COUNT(*) FROM accounts WHERE user_id = $1", input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Account already exists"})
		return
	}

	_, err = h.DB.Exec(`
        INSERT INTO accounts (id, user_id, balance)
        VALUES ($1, $2, 0)
    `, uuid.New(), input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		c.Error(err) // Запишет ошибку в контекст Gin для логов, если настроено
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Account created"})
}
