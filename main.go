// go build && ./api

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Data  []Transaction `json:"data"`
	Count int           `json:"count"`
}

type Transaction struct {
	ID          string `json:"id"`
	AccountID   string `json:"accountId"`
	ToAddress   string `json:"toAddress"`
	FromAddress string `json:"fromAddress"`
	Type        string `json:"type"`   // could also be a custom type or enum for "deposit" or "withdrawal"
	Amount      string `json:"amount"` // string to preserve decimal precision
	Symbol      string `json:"symbol"`
	Decimal     int64  `json:"decimal"`
	Timestamp   string `json:"timestamp"`
	TxnHash     string `json:"txnHash"`
    BlockNumber int64 `json:"blockNumber"`
}

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

	alchemyResponse, err := AlchemyGetTransfers(accountId)
    if err != nil {
        c.JSON(http.StatusNotFound, err) // todo: correct error codes
    }

	response, err := translateAPIResponse(accountId, alchemyResponse)
    if err != nil {
        c.JSON(http.StatusNotFound, err) // todo: correct error codes
    }

	// Respond with the struct marshalled as JSON
	c.JSON(http.StatusOK, response)
}

func translateAPIResponse(accountId string, alchemyResponse *AlchemyAPIResponse) (*APIResponse, error) {
	var response APIResponse
    response.Count = len(alchemyResponse.Result.Transfers)
	response.Data = make([]Transaction, response.Count)
	for i, transfer := range alchemyResponse.Result.Transfers {
		data, err := translateTranferToTransaction(accountId, transfer)
        if err != nil {
            return nil, err
        }
        response.Data[i] = *data
	}
	return &response, nil
}

func translateTranferToTransaction(accountId string, transfer AlchemyTransfer) (*Transaction, error) {
    var t Transaction

	t.ID = transfer.UniqueId
	t.AccountID = accountId
	t.ToAddress = transfer.To
	t.FromAddress = transfer.From
	t.Type = "deposit" // default to deposit, consider self transfer
	if t.AccountID == t.FromAddress {
		t.Type = "withdrawal"
	}
	t.Amount = fmt.Sprint(*transfer.Value) // todo: perhaps only as many digits as decimals in token
	t.Symbol = *transfer.Asset
	t.Timestamp = transfer.Metadata.BlockTimestamp
	t.TxnHash = transfer.Hash

    decimal, err := hexStringToInt(*transfer.RawContract.Decimal)
    if err != nil {
        return nil, err
    }
    t.Decimal = decimal

    blockNumber, err := hexStringToInt(transfer.BlockNum)
    if err != nil {
        return nil, err
    }
    t.BlockNumber = blockNumber

	return &t, nil
}

func hexStringToInt(x string) (int64, error) {
    y, err := strconv.ParseInt(x[2:], 16, 64) // x[2:] to ignore 0x of hex
    if err != nil {
        return 0, err
    }
    return y, nil
}