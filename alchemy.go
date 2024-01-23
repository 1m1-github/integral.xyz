package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Define the struct for the metadata
type AlchemyMetadata struct {
	BlockTimestamp string `json:"blockTimestamp"`
}

// Use a pointer for nullable fields

// Define the struct for the raw contract
type AlchemyRawContract struct {
	Value   *string `json:"value"`
	Address *string `json:"address"`
	Decimal *string `json:"decimal"`
}

// Define the struct for each transfer item
type AlchemyTransfer struct {
	BlockNum        string             `json:"blockNum"`
	UniqueId        string             `json:"uniqueId"`
	Hash            string             `json:"hash"`
	From            string             `json:"from"`
	To              string             `json:"to"`
	Value           *float64           `json:"value"`
	Erc721TokenId   *string            `json:"erc721TokenId"`
	Erc1155Metadata *string            `json:"erc1155Metadata"`
	TokenId         string             `json:"tokenId"`
	Asset           *string            `json:"asset"`
	Category        string             `json:"category"`
	RawContract     AlchemyRawContract `json:"rawContract"`
	Metadata        AlchemyMetadata    `json:"metadata"`
}

// Define the outermost struct
type AlchemyTransfersAPIResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		Transfers []AlchemyTransfer `json:"transfers"`
	} `json:"result"`
}

type AlchemyTokenBalance struct {
	ContractAddress string `json:"contractAddress"`
	TokenBalance    string `json:"tokenBalance"`
}

type AlchemyTokenBalanceResult struct {
	Address       string                `json:"address"`
	TokenBalances []AlchemyTokenBalance `json:"tokenBalances"`
}

type AlchemyBalancesAPIResponse struct {
	JSONRPC string                    `json:"jsonrpc"`
	ID      int                       `json:"id"`
	Result  AlchemyTokenBalanceResult `json:"result"`
}

func AlchemyGetTransfers(accountId string) (*AlchemyTransfersAPIResponse, error) {
	url := "https://eth-mainnet.g.alchemy.com/v2/qmL5zSTAO4Eg3I1O3gdnCommibXbh5Ga"

	payloadStr := fmt.Sprintf(`{"id":1,"jsonrpc":"2.0","method":"alchemy_getAssetTransfers","params":[{"category":["external","internal","erc20","specialnft"],"order":"desc","fromBlock":"0x0","toBlock":"latest","toAddress":"%s","withMetadata":true,"excludeZeroValue":true,"maxCount":"0x3e8"}]}`, accountId)
	payload := strings.NewReader(payloadStr)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response AlchemyTransfersAPIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func AlchemyGetBalances(accountId string) (*AlchemyBalancesAPIResponse, error) {

	url := "https://eth-mainnet.g.alchemy.com/v2/qmL5zSTAO4Eg3I1O3gdnCommibXbh5Ga"

	payloadStr := fmt.Sprintf("{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"alchemy_getTokenBalances\",\"params\":[\"%s\",\"erc20\"]}", accountId)
	payload := strings.NewReader(payloadStr)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	var response AlchemyBalancesAPIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
