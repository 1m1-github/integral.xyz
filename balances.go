package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BalancesAPIResponse struct {
	Data  []Balance `json:"data"`
	Count int       `json:"count"`
}

type Balance struct {
	ContractAddress string `json:"contractAddress"`
	TokenBalance    int64  `json:"tokenBalance"`
}

// BalancesHandler handles requests to the /accounts/:accountId/balances endpoint
func BalancesHandler(c *gin.Context) {
	// Get the accountId from the URL
	accountId := c.Param("accountId")

	alchemyResponse, err := AlchemyGetBalances(accountId)
	if err != nil {
		c.JSON(http.StatusNotFound, err) // todo: correct error codes
	}

	response, err := translateBalancesAPIResponse(accountId, alchemyResponse)
	if err != nil {
		c.JSON(http.StatusNotFound, err) // todo: correct error codes
	}

	// Respond with the struct marshalled as JSON
	c.JSON(http.StatusOK, response)
}

func translateBalancesAPIResponse(accountId string, alchemyResponse *AlchemyBalancesAPIResponse) (*BalancesAPIResponse, error) {
	var response BalancesAPIResponse
	response.Count = len(alchemyResponse.Result.TokenBalances)
	response.Data = make([]Balance, response.Count)
	for i, transfer := range alchemyResponse.Result.TokenBalances {
		data, err := translateBalance(accountId, transfer)
		if err != nil {
			return nil, err
		}
		response.Data[i] = *data
	}
	return &response, nil
}

func translateBalance(accountId string, balance AlchemyTokenBalance) (*Balance, error) {
	var b Balance

	b.ContractAddress = balance.ContractAddress

	balanceValue, err := hexStringToInt(balance.TokenBalance)
	if err != nil {
		return nil, err
	}
	b.TokenBalance = balanceValue

	return &b, nil
}
