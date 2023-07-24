package coingecko

import (
	"encoding/json"
	"fmt"
	"github.com/opal-project-dev/oracle/internal/amount"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func getCurrentPrice(currencyId, conversionCurrency string) (amount.Amount, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s", currencyId), nil)
	if err != nil {
		return amount.Amount{}, err
	}
	req.Header.Set("Accepts", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	bodyBB, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return amount.Amount{}, errors.Errorf("Failed to get current price. Status code: %d, body: %s", resp.StatusCode, string(bodyBB))
	}
	if err != nil {
		return amount.Amount{}, errors.WithMessage(err, "Error sending request to server")
	}
	responseMap := make(map[string]interface{})
	err = json.Unmarshal(bodyBB, &responseMap)
	if err != nil {
		return amount.Amount{}, errors.WithMessage(err, "failed to decode response body.")
	}
	marketData := responseMap["market_data"].(map[string]interface{})
	result := marketData["current_price"].(map[string]interface{})
	am := result[strings.ToLower(conversionCurrency)].(float64)
	return amount.NewFromString(strconv.FormatFloat(am, 'f', -1, 64))
}
