package orca

import (
	"encoding/json"
	"fmt"
)

type AllPools map[string]AllPoolsValue

type AllPoolsValue struct {
	PoolID                string    `json:"poolId"`
	PoolAccount           string    `json:"poolAccount"`
	TokenAAmount          string    `json:"tokenAAmount"`
	TokenBAmount          string    `json:"tokenBAmount"`
	PoolTokenSupply       string    `json:"poolTokenSupply"`
	Apy                   Apy       `json:"apy"`
	Volume                Apy       `json:"volume"`
	RiskLevel             RiskLevel `json:"riskLevel"`
	RiskLevelVolatility   *float64  `json:"riskLevelVolatility,omitempty"`
	RiskLevelModifiedTime *int64    `json:"riskLevelModifiedTime,omitempty"`
}

type Apy struct {
	Day   string `json:"day"`
	Week  string `json:"week"`
	Month string `json:"month"`
}

type RiskLevel string

const (
	Low     RiskLevel = "low"
	Medium  RiskLevel = "medium"
	Minimal RiskLevel = "minimal"
	Unknown RiskLevel = "unknown"
)

func (c *Client) GetAllPools() (AllPools, error) {
	url := fmt.Sprintf("%s/allPools", baseURL)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	var data AllPools
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
