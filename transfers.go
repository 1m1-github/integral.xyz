package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransfersAPIResponse struct {
	Data  []Transfer `json:"data"`
	Count int        `json:"count"`
}

type Transfer struct {
	ID          string `json:"id"`
	AccountID   string `json:"accountId"`
	ToAddress   string `json:"toAddress"`
	FromAddress string `json:"fromAddress"`
	Type        string `json:"type"`
	Amount      string `json:"amount"`
	Symbol      string `json:"symbol"`
	Decimal     int64  `json:"decimal"`
	Timestamp   string `json:"timestamp"`
	TxnHash     string `json:"txnHash"`
	BlockNumber int64  `json:"blockNumber"`
}

// TransactionsHandler handles requests to the /accounts/:accountId/transactions endpoint
func TransactionsHandler(c *gin.Context) {
	// Get the accountId from the URL
	accountId := c.Param("accountId")

	alchemyResponse, err := AlchemyGetTransfers(accountId)
	if err != nil {
		c.JSON(http.StatusNotFound, err) // todo: correct error codes
	}

	response, err := translateTransfersAPIResponse(accountId, alchemyResponse)
	if err != nil {
		c.JSON(http.StatusNotFound, err) // todo: correct error codes
	}

	// Respond with the struct marshalled as JSON
	c.JSON(http.StatusOK, response)
}

func translateTransfersAPIResponse(accountId string, alchemyResponse *AlchemyTransfersAPIResponse) (*TransfersAPIResponse, error) {
	var response TransfersAPIResponse
	response.Count = len(alchemyResponse.Result.Transfers)
	response.Data = make([]Transfer, response.Count)
	for i, transfer := range alchemyResponse.Result.Transfers {
		data, err := translateTranferToTransaction(accountId, transfer)
		if err != nil {
			return nil, err
		}
		response.Data[i] = *data
	}
	return &response, nil
}

func translateTranferToTransaction(accountId string, transfer AlchemyTransfer) (*Transfer, error) {
	var t Transfer

	t.ID = transfer.UniqueId
	t.AccountID = accountId
	t.ToAddress = transfer.To
	t.FromAddress = transfer.From
	t.Type = "deposit" // default to deposit, consider self transfer
	if t.AccountID == t.FromAddress {
		t.Type = "withdrawal"
	}
	t.Amount = fmt.Sprint(*transfer.Value) // todo: perhaps only as many digits as decimals in token, error handling
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
