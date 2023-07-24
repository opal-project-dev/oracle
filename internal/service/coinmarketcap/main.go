package coinmarketcap

import (
	"encoding/json"
	"github.com/opal-project-dev/oracle/internal/amount"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func getCurrentPrice(currencySlug, conversionCurrency, apiKey string) (amount.Amount, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		return amount.Amount{}, err
	}
	addQuery(req, currencySlug, conversionCurrency, apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return amount.Amount{}, errors.WithMessage(err, "Error sending request to server")
	}
	return unmarshalCurrentResponse(resp, conversionCurrency)
}

func addQuery(req *http.Request, currencySlug, conversionCurrency, apiKey string) {
	q := url.Values{}
	q.Add("slug", currencySlug)
	q.Add("convert", conversionCurrency)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	req.URL.RawQuery = q.Encode()
}

func unmarshalCurrentResponse(resp *http.Response, conversionCurrency string) (amount.Amount, error) {
	responseMap := make(map[string]map[string]interface{})
	err := json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return amount.Amount{}, errors.WithMessage(err, "failed to decode response body.")
	}
	currencyId := getKeyOfAnInterface(responseMap["data"])
	currencyKeyValue := responseMap["data"][currencyId]
	quoteKeyValue := currencyKeyValue.(map[string]interface{})
	USD := getValueOfAnInterface(quoteKeyValue["quote"], strings.ToUpper(conversionCurrency))
	price := getValueOfAnInterface(USD, "price").(float64)
	return amount.NewFromString(strconv.FormatFloat(price, 'f', -1, 64))
}

func getKeyOfAnInterface(value interface{}) string {
	c := reflect.ValueOf(value)
	if c.Kind() == reflect.Map {
		return c.MapKeys()[0].Interface().(string)
	}
	return ""
}

func getValueOfAnInterface(value interface{}, key string) interface{} {
	var i interface{}
	c := reflect.ValueOf(value)
	if c.Kind() == reflect.Map {
		for _, k := range c.MapKeys() {
			if k.String() == key {
				strct := c.MapIndex(k)
				i = strct.Interface()
			}
		}
	}
	return i
}
