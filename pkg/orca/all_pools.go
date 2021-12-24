package orca

import (
	"encoding/json"
	"fmt"
)

type Pool struct {
	Name           string   `json:"name"`
	Name2          string   `json:"name2"`
	Account        string   `json:"account"`
	MintAccount    string   `json:"mint_account"`
	Liquidity      float64  `json:"liquidity"`
	Price          float64  `json:"price"`
	Apy24H         *float64 `json:"apy_24h"`
	Apy7D          *float64 `json:"apy_7d"`
	Apy30D         *float64 `json:"apy_30d"`
	Volume24H      float64  `json:"volume_24h"`
	Volume24HQuote float64  `json:"volume_24h_quote"`
	Volume7D       float64  `json:"volume_7d"`
	Volume7DQuote  float64  `json:"volume_7d_quote"`
	Volume30D      float64  `json:"volume_30d"`
	Volume30DQuote float64  `json:"volume_30d_quote"`
}

func (c *Client) GetPools() ([]*Pool, error) {
	url := fmt.Sprintf("%s/pools", baseURL)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	var data []*Pool
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
