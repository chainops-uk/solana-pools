package solana_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/portto/solana-go-sdk/rpc"
)

type GetVoteAccountsResponse struct {
	rpc.GeneralResponse
	Result GetVoteAccountsResult `json:"result"`
}

type GetVoteAccountsResult struct {
	Current []struct {
		Commission       int       `json:"commission"`
		EpochVoteAccount bool      `json:"epochVoteAccount"`
		EpochCredits     [][]int64 `json:"epochCredits"`
		NodePubKey       string    `json:"nodePubkey"`
		LastVote         int       `json:"lastVote"`
		ActivatedStake   int64     `json:"activatedStake"`
		VotePubKey       string    `json:"votePubkey"`
	} `json:"current"`
	Delinquent []struct {
		Commission       int       `json:"commission"`
		EpochVoteAccount bool      `json:"epochVoteAccount"`
		EpochCredits     [][]int64 `json:"epochCredits"`
		NodePubKey       string    `json:"nodePubkey"`
		LastVote         int       `json:"lastVote"`
		ActivatedStake   int64     `json:"activatedStake"`
		VotePubKey       string    `json:"votePubkey"`
	} `json:"delinquent"`
}

func GetVoteAccounts(body []byte, err error) (GetVoteAccountsResult, error) {
	if err != nil {
		return GetVoteAccountsResult{}, fmt.Errorf("rpc: call error, err: %v", err)
	}
	var res GetVoteAccountsResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return GetVoteAccountsResult{}, fmt.Errorf("rpc: failed to json decode body, err: %v", err)
	}
	return res.Result, nil
}
