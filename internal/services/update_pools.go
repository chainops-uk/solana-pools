package services

import (
	"errors"
	"fmt"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/pkg/pools"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

const (
	DefaultTicksPerSecond = 160
	DefaultTicksPerSlot   = 64
	SecondsPerDay         = 60 * 60 * 24
	DefaultSPerSlot       = float64(DefaultTicksPerSlot) / float64(DefaultTicksPerSecond)
	TicksPerDay           = DefaultTicksPerSecond * SecondsPerDay
	DefaultSlotsPerEpoch  = 2 * TicksPerDay / DefaultTicksPerSlot
	SecondsPerEpoch       = DefaultSlotsPerEpoch * DefaultSPerSlot
	EpochsPerYear         = SecondsPerDay * 365.25 / SecondsPerEpoch
)

var (
	nodeAddressNotFounded = errors.New("node address not founded")
)

func (s Imp) UpdatePools() error {
	dPools, err := s.dao.GetPools(nil)
	if err != nil {
		return fmt.Errorf("dao.GetPools: %s", err.Error())
	}
	var success, fail uint64
	start := time.Now()
	for _, p := range dPools {
		if !p.Active {
			continue
		}
		if err := s.updatePool(p); err != nil {
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
	fmt.Println("start: ", dPool.Name)
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

	dmodel := &dmodels.PoolData{
		ID:                uuid.NewV4(),
		PoolID:            dPool.ID,
		ActiveStake:       data.SolanaStake,
		TotalTokensSupply: data.TotalTokenSupply,
		TotalLamports:     data.TotalLamports,
		UnstakeLiquidity:  data.UnstakeLiquidity,
		Epoch:             data.Epoch,
		DepossitFee:       decimal.NewFromFloat(data.DepositFee).Truncate(-2),
		WithdrawalFee:     decimal.NewFromFloat(data.WithdrawalFee).Truncate(-2),
		RewardsFee:        decimal.NewFromFloat(data.RewardsFee).Truncate(-2),
		UpdatedAt:         time.Now(),
		CreatedAt:         time.Now(),
	}

	validatorsPoolData := make([]*dmodels.PoolValidatorData, 0, len(data.Validators))
	for _, v := range data.Validators {

		validator, err := s.dao.GetValidatorByVotePK(v.VotePK)
		if err != nil {
			return fmt.Errorf("dao.GetValidatorByVotePK(%s): %w", v.VotePK, err)
		}
		if validator == nil {
			continue
		}

		validatorsPoolData = append(validatorsPoolData, &dmodels.PoolValidatorData{
			ValidatorID: validator.ID,
			PoolDataID:  dmodel.ID,
			ActiveStake: v.ActiveStake,
		})
	}

	d, err := s.dao.GetLastEpochPoolData(dmodel.PoolID, dmodel.Epoch)
	if err != nil {
		return fmt.Errorf("dao.UpdatePoolData: %w", err)
	}
	if d != nil {
		var epochRate decimal.Decimal
		if d.ActiveStake != 0 {
			lastEpochPoolTokenValue := decimal.NewFromInt(int64(d.TotalLamports)).
				Div(decimal.NewFromInt(int64(d.TotalTokensSupply)))
			TokenValue := decimal.NewFromInt(int64(dmodel.TotalLamports)).
				Div(decimal.NewFromInt(int64(dmodel.TotalTokensSupply)))
			epochRate = TokenValue.Div(lastEpochPoolTokenValue).Sub(decimal.NewFromInt(1))
			epochRate = epochRate.Mul(decimal.NewFromInt(int64(dmodel.Epoch - d.Epoch)))
		} else {
			epochRate = decimal.NewFromInt(0)
		}
		dmodel.APY = epochRate.Mul(decimal.NewFromFloat(EpochsPerYear))
	} else {
		dmodel.APY = decimal.NewFromInt(0)
	}

	err = s.dao.UpdatePoolData(dmodel)
	if err != nil {
		return fmt.Errorf("dao.UpdatePoolData: %s", err.Error())
	}
	err = s.dao.CreatePoolValidatorData(validatorsPoolData...)
	if err != nil {
		return fmt.Errorf("dao.UpdateValidators: %s", err.Error())
	}
	return nil
}
