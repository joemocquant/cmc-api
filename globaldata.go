package coinmarketcap

import (
	"encoding/json"
	"fmt"
)

type GlobalData struct {
	TotalMarketCapUSD            float64 `json:"total_market_cap_usd"`
	Total24hVolumeUSD            float64 `json:"total_24h_volume_usd"`
	BitcoinPercentageOfMarketCap float64 `json:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies             int     `json:"active_currencies"`
	ActiveAsset                  int     `json:"active_assets"`
	ActiveMarkets                int     `json:"active_markets"`
}

// Coinmarketcap API implementation of Gobal Data endpoint
//
// Endpoint: /global/
// Method: GET
//
// Optional parameters:
// (string) convert - return 24h volume, and market cap in terms of another currency. Valid values are:
// "AUD", "BRL", "CAD", "CHF", "CNY", "EUR", "GBP", "HKD", "IDR", "INR", "JPY", "KRW", "MXN", "RUB"
//
// Example: https://api.coinmarketcap.com/v1/global/
// Example: https://api.coinmarketcap.com/v1/global/?convert=EUR
//
// Sample Response:
//
//  {
//    "total_market_cap_usd": 12756692479.0,
//    "total_24h_volume_usd": 135078435.0,
//    "bitcoin_percentage_of_market_cap": 83.34,
//    "active_currencies": 653,
//    "active_assets": 59,
//    "active_markets": 1995
//  }
func (client *Client) GetGlobalData() (*GlobalData, error) {

	resp, err := client.do("global", nil)
	if err != nil {
		return nil, fmt.Errorf("Client.do: %v", err)
	}

	res := GlobalData{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}
