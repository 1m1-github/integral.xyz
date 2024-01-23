// go build && ./api
// todo:
// pagination, input validation, db caching
// integral.xyz accountId, requestId

// testing
// curl http://localhost:8080/accounts/0x9aa99c23f67c81701c772b106b4f83f6e858dd2e/transactions
// curl http://localhost:8080/accounts/0x9aa99c23f67c81701c772b106b4f83f6e858dd2e/balances

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Initialize the Gin engine.
	r := gin.Default()

	// Apply the rate limiter middleware
	// Example: Allow 5 requests per second with a burst of 2
	r.Use(RateLimitMiddleware(5, 2))

	// Define the endpoint with a URL parameter
	r.GET("/accounts/:accountId/transactions", TransactionsHandler)
	r.GET("/accounts/:accountId/balances", BalancesHandler)

	// Start the server on port 8080
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
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