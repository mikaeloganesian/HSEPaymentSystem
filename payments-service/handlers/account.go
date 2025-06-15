package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"payments-service/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	DB *sqlx.DB
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var req struct {
		UserID uuid.UUID `json:"user_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	tx, err := h.DB.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start transaction"})
		return
	}
	defer tx.Rollback() // откат, если что-то пойдёт не так

	// Проверка существующего аккаунта
	var existing models.Account
	err = tx.Get(&existing, "SELECT * FROM accounts WHERE user_id = $1", req.UserID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "account already exists"})
		return
	}

	acc := models.Account{
		UserID:    req.UserID,
		Balance:   0,
		CreatedAt: time.Now(),
	}

	err = tx.QueryRowx(`
		INSERT INTO accounts (user_id, balance, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`, acc.UserID, acc.Balance, acc.CreatedAt).Scan(&acc.ID)
	if err != nil {
		log.Println("Error creating account:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create account"})
		return
	}

	// Добавление в outbox
	err = insertOutboxEvent(tx, "account_created", acc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue event"})
		return
	}

	// Commit транзакции
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, acc)
}

func (h *Handler) Deposit(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := c.BindJSON(&req); err != nil || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	tx, err := h.DB.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start transaction"})
		return
	}
	defer tx.Rollback()

	// Пополнение
	_, err = tx.Exec(`
		UPDATE accounts SET balance = balance + $1 WHERE id = $2
	`, req.Amount, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "deposit failed"})
		return
	}

	// Добавление события в outbox
	payload := map[string]interface{}{
		"account_id": id,
		"amount":     req.Amount,
		"timestamp":  time.Now(),
	}
	err = insertOutboxEvent(tx, "deposit_made", payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue event"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) GetBalance(c *gin.Context) {
	id := c.Param("id")

	var acc models.Account
	err := h.DB.Get(&acc, "SELECT * FROM accounts WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": acc.Balance,
	})
}

// insertOutboxEvent добавляет событие в таблицу outbox в рамках транзакции.
func insertOutboxEvent(tx *sqlx.Tx, eventType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO outbox (event_type, payload, created_at, sent)
		VALUES ($1, $2, $3, false)
	`, eventType, string(data), time.Now())

	return err
}
