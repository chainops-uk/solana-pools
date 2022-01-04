package atrix

import (
	"encoding/json"
	"fmt"
)

type AllPools struct {
	Tvl   float64 `json:"tvl"`
	Pools []Pool  `json:"pools"`
	Farms []Farm  `json:"farms"`
}

type Farm struct {
	Key string  `json:"key"`
	Tvl float64 `json:"tvl"`
	Apy float64 `json:"apy"`
}

type Pool struct {
	PoolKey      string  `json:"poolKey"`
	Tvl          float64 `json:"tvl"`
	LpMint       string  `json:"lpMint"`
	LpSupply     float64 `json:"lpSupply"`
	CoinMint     string  `json:"coinMint"`
	CoinTokens   float64 `json:"coinTokens"`
	CoinDecimals int64   `json:"coinDecimals"`
	PCMint       string  `json:"pcMint"`
	PCTokens     float64 `json:"pcTokens"`
	PCDecimals   int64   `json:"pcDecimals"`
}

func (c *Client) GetTVL() (*AllPools, error) {
	url := fmt.Sprintf("%s/api/tvl", baseURL)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	var data *AllPools
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
