// go build && ./api
// todo:
// pagination, input validation, db caching
// integral.xyz accountId, requestId

// testing
// curl http://localhost:8080/accounts/0x9aa99c23f67c81701c772b106b4f83f6e858dd2e/transactions
// curl http://localhost:8080/accounts/0x9aa99c23f67c81701c772b106b4f83f6e858dd2e/balances

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	fmt.Println("integral.xyz API with endpoints:")
	fmt.Println("/accounts/:accountId/transactions")
	fmt.Println("/accounts/:accountId/balances")

	// Initialize the Gin engine.
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Apply the rate limiter middleware
	// Example: Allow 5 requests per second with a burst of 2
	r.Use(RateLimitMiddleware(5, 2))

	// Define the endpoints with a URL parameter
	r.GET("/accounts/:accountId/transactions", TransactionsHandler)
	r.GET("/accounts/:accountId/balances", BalancesHandler)

	// Start the server on 0.0.0.0:8080
	r.Run(":8080")
}

// RateLimitMiddleware creates a new rate limiter middleware
func RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}