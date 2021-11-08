package solana_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/portto/solana-go-sdk/rpc"
)

type GetInflationRewardResponse struct {
	rpc.GeneralResponse
	Result []GetInflationRewardResult `json:"result"`
}

type GetInflationRewardResult struct {
	Amount        int64 `json:"amount"`
	Commission    int   `json:"commission"`
	EffectiveSlot int   `json:"effectiveSlot"`
	Epoch         int   `json:"epoch"`
	PostBalance   int64 `json:"postBalance"`
}

func GetInflationReward(body []byte, err error) ([]GetInflationRewardResult, error) {
	if err != nil {
		return nil, fmt.Errorf("rpc: call error, err: %v", err)
	}
	var res GetInflationRewardResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("rpc: failed to json decode body, err: %v", err)
	}
	return res.Result, nil
}
