package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/pkg/pools"
	"github.com/everstake/solana-pools/pkg/pools/types"
	"github.com/portto/solana-go-sdk/client"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

func (s Imp) UpdatePools() error {
	dPools, err := s.dao.GetPools()
	if err != nil {
		return fmt.Errorf("dao.GetPools: %s", err.Error())
	}
	var success, fail uint64
	start := time.Now()
	for _, p := range dPools {
		if !p.Active {
			continue
		}
		err = s.updatePool(p)
		if err != nil {
			s.log.Error(
				"Update Pools",
				zap.String("pool_name", p.Name),
				zap.String("pool_address", p.Address),
				zap.String("network", p.Network),
				zap.Error(err),
			)
			fail++
		} else {
			success++
		}
	}
	s.log.Debug(
		"Pools Updated",
		zap.Uint64("success", success),
		zap.Uint64("failed", fail),
		zap.Duration("duration", time.Now().Sub(start)),
	)
	return nil
}

func (s Imp) updatePool(dPool dmodels.Pool) error {
	net := config.Network(dPool.Network)
	rpcCli, ok := s.rpcClients[net]
	if !ok {
		return fmt.Errorf("rpc client for %s network not found", dPool.Network)
	}
	poolFactory := pools.NewFactory(rpcCli)
	pool, err := poolFactory.GetPool(dPool.Name)
	if err != nil {
		return fmt.Errorf("poolFactory.GetPool: %s", err.Error())
	}
	data, err := pool.GetData(dPool.Address)
	if err != nil {
		return fmt.Errorf("pool.GetData: %s", err.Error())
	}
	validatorsMap, err := s.makeValidatorsKeyMap(net)
	if err != nil {
		return fmt.Errorf("makeValidatorsKeyMap: %s", err.Error())
	}
	var validators []dmodels.Validator
	var avgSkippedSlots decimal.Decimal
	var avgScore int64
	var delinquent int64
	for _, v := range data.Validators {
		if v.NodePK == types.EmptyAddress {
			v.NodePK = solana.MustPublicKeyFromBase58(validatorsMap[v.VotePK.String()])
		}
		vInfo, err := s.validatorsApp.GetValidatorInfo(dPool.Network, v.NodePK.String())
		if err != nil {
			return fmt.Errorf("validatorsApp.GetValidatorInfo(%s): %s", v.NodePK, err.Error())
		}
		skippedSlots, _ := decimal.NewFromString(vInfo.SkippedSlotPercent)
		validators = append(validators, dmodels.Validator{
			PoolID: dPool.ID,
			//APR:          0, todo
			VotePK:       v.VotePK.String(),
			NodePK:       v.NodePK.String(),
			ActiveStake:  lampToSol(v.ActiveStake),
			Fee:          decimal.New(vInfo.Commission, 0),
			Score:        vInfo.TotalScore,
			SkippedSlots: skippedSlots,
			DataCenter:   vInfo.DataCenterHost,
		})
		if vInfo.Delinquent {
			delinquent++
		}
		avgSkippedSlots = avgSkippedSlots.Add(skippedSlots)
		avgScore += vInfo.TotalScore
	}
	if len(validators) > 0 {
		avgSkippedSlots = avgSkippedSlots.Div(decimal.New(int64(len(validators)), 0))
		avgScore = avgScore / int64(len(validators))
	}
	err = s.dao.DeleteValidators(dPool.ID)
	if err != nil {
		return fmt.Errorf("dao.DeleteValidators: %s", err.Error())
	}
	err = s.dao.CreateValidators(validators)
	if err != nil {
		return fmt.Errorf("dao.CreateValidators: %s", err.Error())
	}
	dPool.AVGSkippedSlots = avgSkippedSlots
	dPool.AVGScore = avgScore
	dPool.ActiveStake = lampToSol(data.SolanaStake)
	dPool.Nodes = uint64(len(validators))
	if len(validators) > 0 {
		dPool.Delinquent = decimal.NewFromInt(delinquent).Div(decimal.NewFromInt(int64(len(validators))))
	}
	dPool.TokensSupply = decimal.New(int64(data.TokenSupply), -9)
	dPool.DepossitFee = decimal.NewFromFloat(data.DepositFee).Truncate(-2)
	dPool.WithdrawalFee = decimal.NewFromFloat(data.WithdrawalFee).Truncate(-2)
	dPool.RewardsFee = decimal.NewFromFloat(data.RewardsFee).Truncate(-2)
	// todo
	//dPool.UnstakeLiquidity =
	//dPool.APR =
	err = s.dao.UpdatePool(dPool)
	if err != nil {
		return fmt.Errorf("dao.UpdatePool: %s", err.Error())
	}
	return nil
}

// makeValidatorsKeyMap provide map with node and account public keys
func (s Imp) makeValidatorsKeyMap(network config.Network) (mp map[string]string, err error) {
	var cli *client.Client
	switch network {
	case config.Testnet:
		cli = client.NewClient(s.cfg.TestnetNode)
	case config.Mainnet:
		cli = client.NewClient(s.cfg.MainnetNode)
	default:
		return nil, fmt.Errorf("network %s not found", network)
	}
	resp, err := cli.RpcClient.Call(context.Background(), "getVoteAccounts", map[string]string{"commitment": "confirmed"})
	if err != nil {
		return mp, fmt.Errorf("c.Call: %s", err.Error())
	}
	type (
		VotesAccount struct {
			ActivatedStake uint64 `json:"activatedStake"`
			VotePubkey     string `json:"votePubkey"`
			NodePubkey     string `json:"nodePubkey"`
			Commission     uint64 `json:"commission"`
		}
		VotesAccounts struct {
			Result struct {
				Current    []VotesAccount `json:"current"`
				Delinquent []VotesAccount `json:"delinquent"`
			} `json:"result"`
		}
	)
	var voteAccounts VotesAccounts
	err = json.Unmarshal(resp, &voteAccounts)
	if err != nil {
		return mp, fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	mp = make(map[string]string)
	for _, acc := range append(voteAccounts.Result.Current, voteAccounts.Result.Delinquent...) {
		mp[acc.VotePubkey] = acc.NodePubkey
	}
	return mp, nil
}
