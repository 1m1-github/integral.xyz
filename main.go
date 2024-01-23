package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    // Initialize the Gin engine.
    r := gin.Default()

    // Define the endpoint with a URL parameter
    r.GET("/accounts/:accountId/transactions", TransactionsHandler)

    // Start the server on port 8080
    r.Run(":8080") // listen and serve on 0.0.0.0:8080
}

// TransactionsHandler handles requests to the /accounts/:accountId/transactions endpoint
func TransactionsHandler(c *gin.Context) {
    // Get the accountId from the URL
    accountId := c.Param("accountId")

    // For demonstration, we'll just return the accountId in the response
    c.String(http.StatusOK, "Transactions for account ID: %s", accountId)
}
