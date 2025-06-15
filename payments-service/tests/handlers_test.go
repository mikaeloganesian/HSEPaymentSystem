package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"payments-service/db"
	"payments-service/handlers"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// setupTestRouter инициализирует роутер Gin с необходимыми обработчиками
func setupTestRouter() *gin.Engine {
	r := gin.Default()
	db.InitDB()

	h := handlers.Handler{DB: db.DB}

	r.POST("/accounts", h.CreateAccount)
	r.POST("/accounts/:id/deposit", h.Deposit)
	r.GET("/accounts/:id/balance", h.GetBalance)

	return r
}

// helper: создаёт аккаунт и возвращает accountID
func createTestAccount(t *testing.T, e *httpexpect.Expect, userID string) int {
	resp := e.POST("/accounts").
		WithJSON(map[string]string{"user_id": userID}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	resp.ContainsKey("id")
	return int(resp.Value("id").Number().Raw())
}

func TestCreateAccount(t *testing.T) {
	server := httptest.NewServer(setupTestRouter())
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	userID := uuid.New().String()

	resp := e.POST("/accounts").
		WithJSON(map[string]string{"user_id": userID}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	resp.ContainsKey("id")
	resp.ValueEqual("user_id", userID)
}

func TestDeposit(t *testing.T) {
	server := httptest.NewServer(setupTestRouter())
	defer server.Close()

	e := httpexpect.Default(t, server.URL)
	userID := uuid.New().String()
	accountID := createTestAccount(t, e, userID)

	e.POST("/accounts/{id}/deposit").
		WithPath("id", accountID).
		WithJSON(map[string]interface{}{"amount": 100}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("status").
		HasValue("status", "success")
}

func TestGetBalance(t *testing.T) {
	server := httptest.NewServer(setupTestRouter())
	defer server.Close()

	e := httpexpect.Default(t, server.URL)
	userID := uuid.New().String()
	accountID := createTestAccount(t, e, userID)

	// Пополняем баланс, чтобы проверить корректность баланса
	e.POST("/accounts/{id}/deposit").
		WithPath("id", accountID).
		WithJSON(map[string]interface{}{"amount": 150}).
		Expect().
		Status(http.StatusOK)

	e.GET("/accounts/{id}/balance").
		WithPath("id", accountID).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("balance").
		HasValue("balance", 150)
}
