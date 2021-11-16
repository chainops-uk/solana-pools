package v1

import (
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math"
	"net/http"
)

// GetPool godoc
// @Summary get pool
// @Schemes
// @Description get pool
// @Param name path string true "Pool name"
// @Accept json
// @Produce json
// @Success 200 {object} tools.ResponseData{data=PoolDetails} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pool/{name} [get]
func (h *Handler) GetPool(g *gin.Context) (interface{}, error) {
	name := g.Param("name")

	pool, err := h.svc.GetPool(name)
	if err != nil {
		h.log.Error("API GetPoolData", zap.Error(err))
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	return (&PoolDetails{}).Set(&pool), nil
}

// GetPools godoc
// @Summary get pools
// @Schemes
// @Description get pools
// @Accept json
// @Produce json
// @Param offset query number true "offset for aggregation" default(1)
// @Param limit query number true "limit for aggregation" default(10)
// @Param name query string false "stake-pool name"
// @Success 200 {object} tools.ResponseData{data=[]PoolDetails} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pools [get]
func (h *Handler) GetPools(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string `form:"name"`
		Limit  uint64 `form:"limit,default=1"`
		Offset uint64 `form:"offset,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	pools, err := h.svc.GetPools(q.Name, q.Limit, q.Offset)
	if err != nil {
		h.log.Error("API GetPoolData", zap.Error(err))
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	aPools := make([]*PoolDetails, len(pools))
	for i, v := range pools {
		aPools[i] = (&PoolDetails{}).Set(v)
	}

	return aPools, nil
}

// GetTotalPoolsStatistic godoc
// @Summary get statistic
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
		return nil, err
	}

	validators, err := h.svc.GetValidators()
	if err != nil {
		return nil, err
	}

	sc, err := h.svc.GetPoolsStatistic()
	if err != nil {
		return nil, err
	}

	ta, _ := sc.ActiveStake.Float64()
	tu, _ := sc.UnstakeLiquidity.Float64()
	ss, _ := sc.AVGSkippedSlots.Float64()

	APY, _ := apy.Float64()

	usd, err := h.svc.GetPrice()
	if err != nil {
		return nil, err
	}

	USD, _ := usd.Float64()

	return &TotalPoolsStatistic{
		TotalActiveStake:      ta,
		TotalActiveStakePool:  float64(h.svc.GetActiveStake()) * math.Pow(10, -9),
		TotalUnstakeLiquidity: tu,
		TotalValidators:       validators,
		NetworkAPY:            APY,
		Pools:                 poolCount,
		MinPerformanceScore:   sc.MINScore,
		AvgPerformanceScore:   sc.AVGScore,
		MaxPerformanceScore:   sc.MAXScore,
		SkippedSlot:           ss,
		USD:                   USD,
	}, nil
}

type (
	TotalPoolsStatistic struct {
		TotalActiveStakePool  float64 `json:"total_active_stake_pool"`
		TotalActiveStake      float64 `json:"total_active_stake"`
		TotalUnstakeLiquidity float64 `json:"total_unstake_liquidity"`
		TotalValidators       int64   `json:"total_validators"`
		NetworkAPY            float64 `json:"network_apy"`
		Pools                 int64   `json:"pools"`
		MinPerformanceScore   int64   `json:"min_performance_score"`
		AvgPerformanceScore   int64   `json:"avg_performance_score"`
		MaxPerformanceScore   int64   `json:"max_performance_score"`
		SkippedSlot           float64 `json:"skipped_slot"`
		USD                   float64 `json:"usd"`
	}
	Pool struct {
		Address          string          `json:"address"`
		Name             string          `json:"name"`
		ActiveStake      decimal.Decimal `json:"active_stake"`
		TokensSupply     decimal.Decimal `json:"tokens_supply"`
		APR              decimal.Decimal `json:"apr"`
		AVGSkippedSlots  decimal.Decimal `json:"avg_skipped_slots"`
		AVGScore         int64           `json:"avg_score"`
		Delinquent       decimal.Decimal `json:"delinquent"`
		UnstakeLiquidity decimal.Decimal `json:"unstake_liquidity"`
		DepossitFee      decimal.Decimal `json:"depossit_fee"`
		WithdrawalFee    decimal.Decimal `json:"withdrawal_fee"`
		RewardsFee       decimal.Decimal `json:"rewards_fee"`
	}
	PoolDetails struct {
		Pool
		Validators []Validator `json:"validators"`
	}
	Validator struct {
		NodePK       string          `json:"node_pk"`
		APY          decimal.Decimal `json:"apy"`
		VotePK       string          `json:"vote_pk"`
		ActiveStake  decimal.Decimal `json:"active_stake"`
		Fee          decimal.Decimal `json:"fee"`
		Score        int64           `json:"score"`
		SkippedSlots decimal.Decimal `json:"skipped_slots"`
		DataCenter   string          `json:"data_center"`
	}
)

func (pd *PoolDetails) Set(details *smodels.PoolDetails) *PoolDetails {
	pd.Pool.Set(&details.Pool)
	validators := make([]Validator, len(details.Validators))
	for i, validator := range details.Validators {
		validators[i].Set(validator)
	}

	return pd
}

func (pl *Pool) Set(pool *smodels.Pool) *Pool {
	pl.Address = pool.Address
	pl.Name = pool.Address
	pl.ActiveStake = pool.ActiveStake
	pl.TokensSupply = pool.TokensSupply
	pl.APR = pool.APY
	pl.AVGSkippedSlots = pool.AVGSkippedSlots
	pl.AVGScore = pool.AVGScore
	pl.Delinquent = pool.Delinquent
	pl.UnstakeLiquidity = pool.UnstakeLiquidity
	pl.DepossitFee = pool.DepossitFee
	pl.WithdrawalFee = pool.WithdrawalFee
	pl.RewardsFee = pool.RewardsFee

	return pl
}

func (v *Validator) Set(validator *smodels.Validator) *Validator {
	v.NodePK = validator.NodePK
	v.APY = validator.APY
	v.VotePK = validator.VotePK
	v.ActiveStake = validator.ActiveStake
	v.Fee = validator.Fee
	v.Score = validator.Score
	v.SkippedSlots = validator.SkippedSlots
	v.DataCenter = validator.DataCenter

	return v
}
