package v1

import (
	"errors"
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/cache"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math"
	"net/http"
	"time"
)

// GetPool godoc
// @Summary WebSocket
// @Schemes
// @Description get pool
// @Param name path string true "Pool name" default(marinade)
// @Accept json
// @Produce json
// @Success 200 {object} tools.ResponseData{data=pool} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pool/{name} [get]
func (h *Handler) GetPool(g *gin.Context) (interface{}, error) {
	name := g.Param("name")

	resp, err := h.svc.GetPool(name)
	if err != nil {
		h.log.Error("API GetPoolData", zap.Error(err))
		if errors.Is(err, postgres.ErrorRecordNotFounded) {
			return nil, tools.NewStatus(http.StatusBadRequest, fmt.Errorf("%s pool not found", name))
		}

		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	return (&pool{}).Set(&resp.Pool), nil
}

// GetEpoch godoc
// @Summary RestAPI
// @Schemes
// @Description get epoch
// @Accept json
// @Produce json
// @Success 200 {object} tools.ResponseData{data=epoch} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /epoch [get]
func (h *Handler) GetEpoch(ctx *gin.Context) (interface{}, error) {
	e, err := h.svc.GetEpoch()
	if err != nil {
		return nil, err
	}

	return tools.ResponseData{Data: (&epoch{}).Set(e)}, nil
}

// GetPools godoc
// @Summary RestAPI
// @Schemes
// @Description get pools
// @Accept json
// @Produce json
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Param name query string false "stake-pool name"
// @Success 200 {object} tools.ResponseArrayData{data=[]poolMainPage} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pools [get]
func (h *Handler) GetPools(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string `form:"name"`
		Offset uint64 `form:"offset,default=0"`
		Limit  uint64 `form:"limit,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	pools, amounnt, err := h.svc.GetPools(q.Name, q.Limit, q.Offset)
	if err != nil {
		h.log.Error("API GetPoolData", zap.Error(err))
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	aPools := make([]*poolMainPage, len(pools))
	for i, v := range pools {
		aPools[i] = (&poolMainPage{}).Set(v)
	}

	return tools.ResponseArrayData{
		Data: aPools,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: amounnt,
		}}, nil
}

// GetTotalPoolsStatistic godoc
// @Summary WebSocket
// @Schemes
// @Description get statistic
// @Accept json
// @Produce json
// @Success 200 {object} tools.ResponseData{data=TotalPoolsStatistic} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pools-statistic [get]
func (h *Handler) GetTotalPoolsStatistic(ctx *gin.Context) (interface{}, error) {
	poolCount, err := h.svc.GetPoolCount()
	if err != nil {
		return nil, err
	}

	apy, err := h.svc.GetAPY()
	if err != nil {
		if errors.Is(err, cache.KeyWasNotFound) {
			return nil, tools.NewStatus(500, fmt.Errorf("apy metric is empty"))
		}
		return nil, err
	}

	validators, err := h.svc.GetValidators()
	if err != nil {
		if errors.Is(err, cache.KeyWasNotFound) {
			return nil, tools.NewStatus(500, fmt.Errorf("validator metric is empty"))
		}
		return nil, err
	}

	sc, err := h.svc.GetPoolsCurrentStatistic()
	if err != nil {
		return nil, err
	}

	ta, _ := sc.ActiveStake.Float64()
	ts, _ := sc.TotalSupply.Float64()
	tu, _ := sc.UnstakeLiquidity.Float64()
	ss, _ := sc.AVGSkippedSlots.Float64()
	paa, _ := sc.AVGPoolsApy.Float64()

	APY, _ := apy.Float64()

	usd, err := h.svc.GetPrice()
	if err != nil {
		return nil, err
	}

	USD, _ := usd.Float64()

	return tools.ResponseData{Data: &TotalPoolsStatistic{
		TotalActiveStake:      float64(h.svc.GetActiveStake()) * math.Pow(10, -9),
		TotalActiveStakePool:  ta,
		TotalSupply:           ts,
		TotalUnstakeLiquidity: tu,
		TotalValidators:       validators,
		NetworkAPY:            APY,
		Pools:                 poolCount,
		PoolsAvgAPY:           paa,
		MinPerformanceScore:   sc.MINScore,
		AvgPerformanceScore:   sc.AVGScore,
		MaxPerformanceScore:   sc.MAXScore,
		SkippedSlot:           ss,
		USD:                   USD,
	}}, nil
}

// GetPoolsStatistic godoc
// @Summary RestAPI
// @Schemes
// @Description get statistic by pool
// @Accept json
// @Produce json
// @Param name query string true "pool name" default(mSOL)
// @Param aggregation query string true "aggregation" Enums(week, month, year)
// @Success 200 {object} tools.ResponseData{data=[]poolStatistic} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pool-statistic [get]
func (h *Handler) GetPoolsStatistic(ctx *gin.Context) (interface{}, error) {
	request := struct {
		Name        string `form:"name" binding:"required"`
		Aggregation string `form:"aggregation" binding:"required"`
	}{}

	if err := ctx.ShouldBind(&request); err != nil {
		return nil, tools.NewStatus(http.StatusNotAcceptable, fmt.Errorf("bad request %w", err))
	}

	arr, err := h.svc.GetPoolStatistic(request.Name, request.Aggregation)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFounded) {
			return nil, tools.NewStatus(http.StatusBadRequest, fmt.Errorf("%s pool not found", request.Name))
		}
		return nil, err
	}

	data := make([]*poolStatistic, len(arr))
	for i, v := range arr {
		data[i] = (&poolStatistic{}).Set(v)
	}

	return tools.ResponseData{Data: data}, err
}

type (
	epoch struct {
		Epoch        uint64    `json:"epoch"`
		SlotsInEpoch uint64    `json:"slots_in_epoch"`
		SPS          float64   `json:"sps"`
		EndEpoch     time.Time `json:"end_epoch"`
		Progress     uint8     `json:"progress"`
	}
	poolMainPage struct {
		pool
		Validators uint64 `json:"validators"`
	}
	poolStatistic struct {
		ActiveStake        float64   `json:"active_stake"`
		APY                float64   `json:"apy"`
		UnstackedLiquidity float64   `json:"unstacked_liquidity"`
		NumberOfValidators int64     `json:"number_of_validators"`
		CreatedAt          time.Time `json:"created_at"`
	}
	TotalPoolsStatistic struct {
		TotalActiveStakePool  float64 `json:"total_active_stake_pool"`
		TotalActiveStake      float64 `json:"total_active_stake"`
		TotalSupply           float64 `json:"total_supply"`
		TotalUnstakeLiquidity float64 `json:"total_unstake_liquidity"`
		TotalValidators       int64   `json:"total_validators"`
		NetworkAPY            float64 `json:"network_apy"`
		Pools                 int64   `json:"pools"`
		PoolsAvgAPY           float64 `json:"pools_avg_apy"`
		MinPerformanceScore   int64   `json:"min_performance_score"`
		AvgPerformanceScore   int64   `json:"avg_performance_score"`
		MaxPerformanceScore   int64   `json:"max_performance_score"`
		SkippedSlot           float64 `json:"skipped_slot"`
		USD                   float64 `json:"usd"`
	}
	pool struct {
		Address          string  `json:"address"`
		Name             string  `json:"name"`
		ThumbImage       string  `json:"thumb_image"`
		SmallImage       string  `json:"small_image"`
		LargeImage       string  `json:"large_image"`
		Currency         string  `json:"currency"`
		ActiveStake      float64 `json:"active_stake"`
		TokensSupply     float64 `json:"tokens_supply"`
		APY              float64 `json:"apy"`
		AVGSkippedSlots  float64 `json:"avg_skipped_slots"`
		AVGScore         int64   `json:"avg_score"`
		StakingAccounts  uint64  `json:"staking_accounts"`
		Delinquent       float64 `json:"delinquent"`
		UnstakeLiquidity float64 `json:"unstake_liquidity"`
		DepossitFee      float64 `json:"depossit_fee"`
		WithdrawalFee    float64 `json:"withdrawal_fee"`
		RewardsFee       float64 `json:"rewards_fee"`
	}
)

func (e *epoch) Set(data *smodels.EpochInfo) *epoch {
	e.Epoch = data.Epoch
	e.SlotsInEpoch = data.SlotsInEpoch
	e.SPS = data.SPS
	e.EndEpoch = data.EndEpoch
	e.Progress = data.Progress
	return e
}

func (ps *poolStatistic) Set(data *smodels.Pool) *poolStatistic {
	ps.APY, _ = data.APY.Float64()
	ps.UnstackedLiquidity, _ = data.UnstakeLiquidity.Float64()
	ps.ActiveStake, _ = data.ActiveStake.Float64()
	ps.NumberOfValidators = data.ValidatorCount
	ps.CreatedAt = data.CreatedAt
	return ps
}

func (pd *poolMainPage) Set(details *smodels.PoolDetails) *poolMainPage {
	pd.pool.Set(&details.Pool)

	return pd
}

func (pl *pool) Set(pool *smodels.Pool) *pool {
	pl.Address = pool.Address
	pl.Name = pool.Name
	pl.SmallImage = pool.SmallImage
	pl.LargeImage = pool.LargeImage
	pl.ThumbImage = pool.ThumbImage
	pl.Currency = pool.Currency
	pl.ActiveStake, _ = pool.ActiveStake.Float64()
	pl.TokensSupply, _ = pool.TokensSupply.Float64()
	pl.APY, _ = pool.APY.Float64()
	pl.AVGSkippedSlots, _ = pool.AVGSkippedSlots.Float64()
	pl.AVGScore = pool.AVGScore
	pl.StakingAccounts = pool.StakingAccounts
	pl.Delinquent, _ = pool.Delinquent.Float64()
	pl.UnstakeLiquidity, _ = pool.UnstakeLiquidity.Float64()
	pl.DepossitFee, _ = pool.DepossitFee.Float64()
	pl.WithdrawalFee, _ = pool.WithdrawalFee.Float64()
	pl.RewardsFee, _ = pool.RewardsFee.Float64()

	return pl
}
