package solana_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/portto/solana-go-sdk/rpc"
)

type GetFeesResponse struct {
	rpc.GeneralResponse
	Result GetFeesResult `json:"result"`
}

type GetFeesResult struct {
	Context rpc.Context  `json:"context"`
	Value   GetFeesValue `json:"value"`
}

type GetFeesValue struct {
	Blockhash     string `json:"blockhash"`
	FeeCalculator struct {
		LamportsPerSignature int `json:"lamportsPerSignature"`
	} `json:"feeCalculator"`
	LastValidBlockHeight int `json:"lastValidBlockHeight"`
	LastValidSlot        int `json:"lastValidSlot"`
}

func GetFees(body []byte, err error) (GetFeesResult, error) {
	if err != nil {
		return GetFeesResult{}, fmt.Errorf("rpc: call error, err: %v", err)
	}
	var res GetFeesResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return GetFeesResult{}, fmt.Errorf("rpc: failed to json decode body, err: %v", err)
	}
	return res.Result, nil
}
