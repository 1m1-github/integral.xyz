package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

    alchemyGetTransfers(accountId)
}

func alchemyGetTransfers(accountId string) {
	url := "https://eth-mainnet.g.alchemy.com/v2/qmL5zSTAO4Eg3I1O3gdnCommibXbh5Ga"

	payloadStr := fmt.Sprintf(`{"id":1,"jsonrpc":"2.0","method":"alchemy_getAssetTransfers","params":[{"category":["external","internal","erc20","specialnft"],"order":"desc","fromBlock":"0x0","toBlock":"latest","toAddress":"%s","withMetadata":true,"excludeZeroValue":true,"maxCount":"0x3e8"}]}`, accountId)
	payload := strings.NewReader(payloadStr)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

}
