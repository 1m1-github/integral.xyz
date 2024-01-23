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
	TokenId         string            `json:"tokenId"`
	Asset           *string             `json:"asset"`
	Category        string             `json:"category"`
	RawContract     AlchemyRawContract `json:"rawContract"`
	Metadata        AlchemyMetadata    `json:"metadata"`
}

// Define the outermost struct
type AlchemyAPIResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		Transfers []AlchemyTransfer `json:"transfers"`
	} `json:"result"`
}

func AlchemyGetTransfers(accountId string) AlchemyAPIResponse {
	url := "https://eth-mainnet.g.alchemy.com/v2/qmL5zSTAO4Eg3I1O3gdnCommibXbh5Ga"

	payloadStr := fmt.Sprintf(`{"id":1,"jsonrpc":"2.0","method":"alchemy_getAssetTransfers","params":[{"category":["external","internal","erc20","specialnft"],"order":"desc","fromBlock":"0x0","toBlock":"latest","toAddress":"%s","withMetadata":true,"excludeZeroValue":true,"maxCount":"0x3e8"}]}`, accountId)
	payload := strings.NewReader(payloadStr)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response AlchemyAPIResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	return response
}
