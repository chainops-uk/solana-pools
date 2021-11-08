package solana_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/portto/solana-go-sdk/rpc"
)

type GetSupplyInfoResponse struct {
	rpc.GeneralResponse
	Result GetSupplyInfoResult `json:"result"`
}

type GetSupplyInfoResult struct {
	Context rpc.Context          `json:"context"`
	Value   GetSupplyResultValue `json:"value"`
}

type GetSupplyResultValue struct {
	Circulating            int64    `json:"circulating"`
	NonCirculating         int64    `json:"nonCirculating"`
	NonCirculatingAccounts []string `json:"nonCirculatingAccounts"`
	Total                  int64    `json:"total"`
}

func GetSupply(body []byte, err error) (GetSupplyInfoResult, error) {
	if err != nil {
		return GetSupplyInfoResult{}, fmt.Errorf("rpc: call error, err: %v", err)
	}
	var res GetSupplyInfoResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return GetSupplyInfoResult{}, fmt.Errorf("rpc: failed to json decode body, err: %v", err)
	}
	return res.Result, nil
}
