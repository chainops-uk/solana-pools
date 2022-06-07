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
// @Description Creates a WS request to get server data for the pool with the pool name specified in the request.
// @Tags pool
// @Param name path string true "Name of the pool with strict observance of the case." default(Eversol)
// @Param epoch query number false "Epoch aggregation." Enums(1, 10) default(10)
// @Accept json
// @Produce json
// @Success 200 {object} tools.ResponseData{data=pool} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pool/{name} [get]
func (h *Handler) GetPool(ctx *gin.Context, message []byte) (interface{}, error) {
	name := ctx.Param("name")

	q := struct {
		Epoch uint64 `form:"epoch,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	resp, err := h.svc.GetPool(name, q.Epoch)
	if err != nil {
		h.log.Error("API GetPoolData", zap.Error(err))
		if errors.Is(err, postgres.ErrorRecordNotFounded) {
			return nil, tools.NewStatus(http.StatusBadRequest, fmt.Errorf("%s pool not found", name))
		}

		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	return (&pool{}).Set(&resp.Pool), nil
}

// GetPools godoc
// @Summary RestAPI
// @Schemes
// @Description This Pools list with ability to sort & search by name.
// @Tags pool
// @Accept json
// @Produce json
// @Param name query string false "The name of the pool without strict observance of the case."
// @Param epoch query number true "Epoch aggregation." Enums(1, 10) default(10)
// @Param sort query string false "The parameter by the value of which the pools will be sorted." Enums(apy, pool stake, validators, score, skipped slot, token price) default(apy)
// @Param desc query bool false "Sort in descending order" default(true)
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Success 200 {object} tools.ResponseArrayData{data=[]poolMainPage} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pools [get]
func (h *Handler) GetPools(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string `form:"name"`
		Epoch  uint64 `form:"epoch,default=10"`
		Sort   string `form:"sort,default=apy"`
		Desc   bool   `form:"desc,default=true"`
		Offset uint64 `form:"offset,default=0"`
		Limit  uint64 `form:"limit,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	pools, amount, err := h.svc.GetPools(q.Name, q.Sort, q.Desc, q.Epoch, q.Limit, q.Offset)
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
			TotalAmount: amount,
		}}, nil
}

// GetTotalPoolsStatistic godoc
// @Summary WebSocket
// @Schemes
// @Description Creates a WS request to get current statistics.
// @Tags pool
// @Param epoch query number false "Epoch aggregation." Enums(1, 10) default(10)
// @Accept json
// @Produce json
// @Success 200 {object} tools.ResponseData{data=TotalPoolsStatistic} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pools-statistic [get]
func (h *Handler) GetTotalPoolsStatistic(ctx *gin.Context, message []byte) (interface{}, error) {
	q := struct {
		Epoch uint64 `form:"epoch,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
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
			return nil, tools.NewStatus(500, fmt.Errorf("validatorData metric is empty"))
		}
		return nil, err
	}

	sc, err := h.svc.GetPoolsCurrentStatistic(q.Epoch)
	if err != nil {
		return nil, err
	}

	ta, _ := sc.ActiveStake.Float64()
	ts, _ := sc.TotalSupply.Float64()
	tu, _ := sc.UnstakeLiquidity.Float64()
	ss, _ := sc.AVGSkippedSlots.Float64()
	paa, _ := sc.MAXPoolsApy.Float64()

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
		Pools:                 sc.Pools,
		PoolsMaxAPY:           paa,
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
// @Description The pool statistic for the specified aggregation.
// @Tags pool
// @Accept json
// @Produce json
// @Param name query string true "Name of the pool with strict observance of the case." default(Eversol)
// @Param aggregation query string true "Type of data aggregation for a time period" Enums(week, month, quarter, half-year, year)
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
		TotalSol           float64   `json:"total_sol"`
		TokensSupply       float64   `json:"tokens_supply"`
		ActiveStake        float64   `json:"active_stake"`
		APY                float64   `json:"apy"`
		Delinquent         uint64    `json:"delinquent"`
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
		Pools                 uint64  `json:"pools"`
		PoolsMaxAPY           float64 `json:"pools_max_apy"`
		MinPerformanceScore   int64   `json:"min_performance_score"`
		AvgPerformanceScore   int64   `json:"avg_performance_score"`
		MaxPerformanceScore   int64   `json:"max_performance_score"`
		SkippedSlot           float64 `json:"skipped_slot"`
		USD                   float64 `json:"usd"`
	}
	pool struct {
		Address          string  `json:"address"`
		Name             string  `json:"name"`
		Image            string  `json:"image"`
		Currency         string  `json:"currency"`
		ActiveStake      float64 `json:"active_stake"`
		TokensSupply     float64 `json:"tokens_supply"`
		TotalSol         float64 `json:"total_sol"`
		APY              float64 `json:"apy"`
		Validators       int64   `json:"validators"`
		AVGSkippedSlots  float64 `json:"avg_skipped_slots"`
		AVGScore         int64   `json:"avg_score"`
		StakingAccounts  uint64  `json:"staking_accounts"`
		Delinquent       uint64  `json:"delinquent"`
		UnstakeLiquidity float64 `json:"unstake_liquidity"`
		DepositFee       float64 `json:"deposit_fee"`
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
	ps.TotalSol, _ = data.TotalLamports.Float64()
	ps.TokensSupply, _ = data.TokensSupply.Float64()
	ps.APY, _ = data.APY.Float64()
	ps.UnstackedLiquidity, _ = data.UnstakeLiquidity.Float64()
	ps.ActiveStake, _ = data.ActiveStake.Float64()
	ps.NumberOfValidators = data.ValidatorCount
	ps.Delinquent = data.Delinquent
	ps.CreatedAt = data.CreatedAt
	return ps
}

func (pd *poolMainPage) Set(details *smodels.PoolDetails) *poolMainPage {
	pd.pool.Set(&details.Pool)
	pd.Validators = uint64(details.ValidatorCount)
	return pd
}

func (pl *pool) Set(pool *smodels.Pool) *pool {
	pl.Address = pool.Address
	pl.Name = pool.Name
	pl.Image = pool.Image
	pl.Currency = pool.Currency
	pl.ActiveStake, _ = pool.ActiveStake.Float64()
	pl.TokensSupply, _ = pool.TokensSupply.Float64()
	pl.TotalSol, _ = pool.TotalLamports.Float64()
	pl.APY, _ = pool.APY.Float64()
	pl.AVGSkippedSlots, _ = pool.AVGSkippedSlots.Float64()
	pl.AVGScore = pool.AVGScore
	pl.StakingAccounts = pool.StakingAccounts
	pl.Delinquent = pool.Delinquent
	pl.UnstakeLiquidity, _ = pool.UnstakeLiquidity.Float64()
	pl.DepositFee, _ = pool.DepossitFee.Float64()
	pl.WithdrawalFee, _ = pool.WithdrawalFee.Float64()
	pl.RewardsFee, _ = pool.RewardsFee.Float64()
	pl.Validators = pool.ValidatorCount

	return pl
}
