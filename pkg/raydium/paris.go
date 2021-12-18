package raydium

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Pairs struct {
	Name            string   `json:"name"`
	PairID          string   `json:"pair_id"`
	LpMint          string   `json:"lp_mint"`
	Official        bool     `json:"official"`
	Liquidity       float64  `json:"liquidity"`
	Market          string   `json:"market"`
	Volume24H       float64  `json:"volume_24h"`
	Volume24HQuote  float64  `json:"volume_24h_quote"`
	Fee24H          float64  `json:"fee_24h"`
	Fee24HQuote     float64  `json:"fee_24h_quote"`
	Volume7D        float64  `json:"volume_7d"`
	Volume7DQuote   float64  `json:"volume_7d_quote"`
	Fee7D           float64  `json:"fee_7d"`
	Fee7DQuote      float64  `json:"fee_7d_quote"`
	Price           *float64 `json:"price"`
	LpPrice         *float64 `json:"lp_price"`
	AmmID           string   `json:"amm_id"`
	TokenAmountCoin float64  `json:"token_amount_coin"`
	TokenAmountPC   float64  `json:"token_amount_pc"`
	TokenAmountLp   float64  `json:"token_amount_lp"`
	Apy             float64  `json:"apy"`
}

func (c *Client) GetPairs(amm_id string) ([]*Pairs, error) {
	params := url.Values{}
	if amm_id != "" {
		params.Add("amm_id", amm_id)
	}

	url := fmt.Sprintf("%s/pairs?%s", baseURL, params.Encode())
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	var data []*Pairs
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
