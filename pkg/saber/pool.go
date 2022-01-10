package saber

import (
	"encoding/json"
)

type AllPools struct {
	Data Data `json:"data"`
}

type Data struct {
	Pools []*Pool `json:"pools"`
}

type Pool struct {
	AmmID string `json:"ammId"`
	Name  string `json:"name"`
	Coin  Coin   `json:"coin"`
	PC    Coin   `json:"pc"`
	Lp    Coin   `json:"lp"`
	Stats Stats  `json:"stats"`
}

type Coin struct {
	ChainID  *int64  `json:"chainId"`
	Address  string  `json:"address"`
	Name     string  `json:"name"`
	Decimals int64   `json:"decimals"`
	Symbol   string  `json:"symbol"`
	LogoURI  *string `json:"logoURI"`
}

type Stats struct {
	TvlPC   float64 `json:"tvl_pc"`
	TvlCoin float64 `json:"tvl_coin"`
	Price   float64 `json:"price"`
	Vol24H  float64 `json:"vol24h"`
}

func (c *Client) GetPools() ([]*Pool, error) {
	body := `{"query":"query AllPoolStats {pools {ammId   name    coin {     chainId     address      name     decimals     symbol      logoURI }   pc {     chainId     address     name    decimals     symbol      logoURI   }   lp {     chainId    address    name     decimals     symbol     logoURI  }    stats {     tvl_pc     tvl_coin      price     vol24h }}}","operationName":"AllPoolStats"}`

	resp, err := c.MakeReq(baseURL, []byte(body))
	if err != nil {
		return nil, err
	}

	var data AllPools
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data.Data.Pools, nil
}
