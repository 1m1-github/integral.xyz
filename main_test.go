package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTransactionsHandler(t *testing.T) {
	router := gin.Default()
	router.GET("/accounts/:accountId/transactions", TransactionsHandler)

	req, _ := http.NewRequest("GET", "/accounts/0x9aa99c23f67c81701c772b106b4f83f6e858dd2e/transactions", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// generic http check
	assert.Equal(t, http.StatusOK, resp.Code)

	var response TransfersAPIResponse
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 0x9aa99c23f67c81701c772b106b4f83f6e858dd2e has 8 transfers
	assert.Equal(t, 8, response.Count)
}

func TestBalancesHandler(t *testing.T) {
	router := gin.Default()
	router.GET("/accounts/:accountId/balances", BalancesHandler)

	req, _ := http.NewRequest("GET", "/accounts/0x9aa99c23f67c81701c772b106b4f83f6e858dd2e/balances", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// generic http check
	assert.Equal(t, http.StatusOK, resp.Code)

	var response BalancesAPIResponse
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 0x9aa99c23f67c81701c772b106b4f83f6e858dd2e has 2 balances
	assert.Equal(t, 2, response.Count)
}
