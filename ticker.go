package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Ticks []*Tick

type Tick struct {
	Id               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int     `json:"rank,string"`
	PriceUSD         float64 `json:"price_usd,string"`
	PriceBTC         float64 `json:"price_btc,string"`
	DayVolumeUSD     float64 `json:"24h_volume_usd,string"`
	MarketCapUSD     float64 `json:"market_cap_usd,string"`
	AvailableSupply  float64 `json:"available_supply,string"`
	TotalSupply      float64 `json:"total_supply,string"`
	PercentChange1h  float64 `json:"percent_change_1h,string"`
	PercentChange24h float64 `json:"percent_change_24h,string"`
	PercentChange7d  float64 `json:"percent_change_7d,string"`
	LastUpdated      int64   `json:"last_updated,string"` // Unix timestamp
}

// Coinmarketcap API implementation of Ticker endpoint.
//
// Endpoint: /ticker/
// Method: GET
//
// Optional parameters:
// (int) limit - only returns the top limit results.
// (string) convert - return price, 24h volume, and market cap in terms of another currency.
// Valid values are: "AUD", "BRL", "CAD", "CHF", "CNY", "EUR", "GBP", "HKD", "IDR", "INR",
// "JPY", "KRW", "MXN", "RUB"
//
// Example: https://api.coinmarketcap.com/v1/ticker/
// Example: https://api.coinmarketcap.com/v1/ticker/?limit=10
// Example: https://api.coinmarketcap.com/v1/ticker/?convert=EUR&limit=10
//
// Sample Response:
//
//  [
//    {
//      "id": "bitcoin",
//      "name": "Bitcoin",
//      "symbol": "BTC",
//      "rank": "1",
//      "price_usd": "573.137",
//      "price_btc": "1.0",
//      "24h_volume_usd": "72855700.0",
//      "market_cap_usd": "9080883500.0",
//      "available_supply": "15844176.0",
//      "total_supply": "15844176.0",
//      "percent_change_1h": "0.04",
//      "percent_change_24h": "-0.3",
//      "percent_change_7d": "-0.57",
//      "last_updated": "1472762067"
//    },
//    {
//      "id": "ethereum",
//      "name": "Ethereum",
//      "symbol": "ETH",
//      "rank": "2",
//      "price_usd": "12.1844",
//      "price_btc": "0.021262",
//      "24h_volume_usd": "24085900.0",
//      "market_cap_usd": "1018098455.0",
//      "available_supply": "83557537.0",
//      "total_supply": "83557537.0",
//      "percent_change_1h": "-0.58",
//      "percent_change_24h": "6.34",
//      "percent_change_7d": "8.59",
//      "last_updated": "1472762062"
//    },
//    ...
//  ]
func (client *Client) GetTickers() (Ticks, error) {

	resp, err := client.do("ticker", nil)
	if err != nil {
		return nil, fmt.Errorf("Client.do: %v", err)
	}

	res := make(Ticks, 0)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

func (client *Client) GetTickersLimit(limit int) (Ticks, error) {

	params := map[string]string{
		"limit": strconv.Itoa(limit),
	}

	resp, err := client.do("ticker", params)
	if err != nil {
		return nil, fmt.Errorf("Client.do: %v", err)
	}

	res := make(Ticks, 0)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

func (client *Client) GetTicker(id string) (*Tick, error) {

	resp, err := client.do("ticker"+"/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("Client.do: %v", err)
	}

	res := make(Ticks, 1)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res[0], nil
}
