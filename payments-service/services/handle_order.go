package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"payments-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (w *Worker) handleOrder(event models.OrderEvent) error {
	tx, err := w.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем: уже обрабатывали это сообщение?
	var exists bool
	err = tx.Get(&exists, "SELECT EXISTS(SELECT 1 FROM inbox WHERE id = $1)", event.ID)
	if err != nil {
		return err
	}
	if exists {
		log.Println("event already processed:", event.ID)
		return nil
	}

	// Добавляем в inbox

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON: %w", err)
	}

	_, err = tx.Exec(`
    INSERT INTO inbox (id, event_type, payload, processed_at)
    VALUES ($1, $2, $3, $4)
`, event.ID, "OrderCreated", payload, time.Now())
	if err != nil {
		return err
	}

	// Пробуем списать деньги
	var balance int64
	err = tx.Get(&balance, "SELECT balance FROM accounts WHERE user_id = $1", event.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No account found for user_id=%s", event.UserID)
			return fmt.Errorf("account not found for user_id %s", event.UserID)
		}
		return err
	}

	if balance < event.Amount {
		log.Println("insufficient balance")
		// outbox: payment_failed
		insertOutbox(tx, "PaymentFailed", event)
		return tx.Commit()
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE user_id = $2", event.Amount, event.UserID)
	if err != nil {
		return err
	}

	// outbox: payment_success
	err = insertOutbox(tx, "PaymentSuccess", event)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func insertOutbox(tx *sqlx.Tx, eventType string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        INSERT INTO outbox (id, event_type, payload, sent, created_at)
        VALUES ($1, $2, $3, false, now())
    `, uuid.New(), eventType, payload)
	return err
}
