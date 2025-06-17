package main

import (
	_ "gateway/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// @title        Gateway API
// @version      1.0
// @description  Gateway между Orders и Payments микросервисами
// @host         localhost:3000
// @BasePath     /api

const (
	paymentsURL = "http://localhost:8080"
	ordersURL   = "http://localhost:8081"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func proxyGET(c *gin.Context, targetURL string) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.URL.RawQuery = c.Request.URL.RawQuery

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func proxyPOST(c *gin.Context, targetURL string) {
	req, err := http.NewRequest("POST", targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header = c.Request.Header

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// getSummary godoc
// @Summary Получить сводную информацию (баланс и заказы)
// @Description Возвращает баланс пользователя и его заказы
// @Tags Summary
// @Param user_id query string true "User ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /summary [get]
func getSummary(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	balanceResp, err := httpClient.Get(fmt.Sprintf("%s/accounts/balance?user_id=%s", paymentsURL, userID))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "payments service unavailable"})
		return
	}
	defer balanceResp.Body.Close()
	balanceBody, _ := io.ReadAll(balanceResp.Body)

	ordersResp, err := httpClient.Get(fmt.Sprintf("%s/orders", ordersURL))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "orders service unavailable"})
		return
	}
	defer ordersResp.Body.Close()
	ordersBody, _ := io.ReadAll(ordersResp.Body)

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"balance": gin.H{"raw": string(balanceBody)},
		"orders":  gin.H{"raw": string(ordersBody)},
	})
}

func main() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.POST("/accounts", func(c *gin.Context) {
			proxyPOST(c, paymentsURL+"/accounts")
		})
		api.POST("/accounts/:id/deposit", func(c *gin.Context) {
			proxyPOST(c, fmt.Sprintf("%s/accounts/%s/deposit", paymentsURL, c.Param("id")))
		})
		api.GET("/accounts/balance", func(c *gin.Context) {
			proxyGET(c, paymentsURL+"/accounts/balance")
		})

		api.POST("/orders", func(c *gin.Context) {
			proxyPOST(c, ordersURL+"/orders")
		})
		api.GET("/order/:id", func(c *gin.Context) {
			proxyGET(c, fmt.Sprintf("%s/order/%s", ordersURL, c.Param("id")))
		})
		api.GET("/orders", func(c *gin.Context) {
			proxyGET(c, ordersURL+"/orders")
		})

		api.GET("/summary", getSummary)
	}

	log.Println("Gateway running on http://localhost:3000")
	router.Run(":3000")
}
