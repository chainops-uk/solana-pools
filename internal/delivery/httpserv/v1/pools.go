package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"net/http"
)

type (
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
		APR          decimal.Decimal `json:"apr"`
		VotePK       string          `json:"vote_pk"`
		ActiveStake  decimal.Decimal `json:"active_stake"`
		Fee          decimal.Decimal `json:"fee"`
		Score        int64           `json:"score"`
		SkippedSlots decimal.Decimal `json:"skipped_slots"`
		DataCenter   string          `json:"data_center"`
	}
)

// GetPool godoc
// @Summary get pool
// @Schemes
// @Description get pool
// @Param name path string true "Pool name"
// @Accept json
// @Produce json
// @Success 200 {object} PoolDetails
// @Router /pool/{name} [get]
func (h Handler) GetPool(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		g.JSON(http.StatusBadRequest, nil)
		return
	}
	pool, err := h.svc.GetPool(name)
	if err != nil {
		h.log.Error("API GetPoolData", zap.Error(err))
		g.JSON(http.StatusInternalServerError, nil)
		return
	}
	validators := make([]Validator, len(pool.Validators))
	for i, v := range pool.Validators {
		validators[i] = Validator{
			NodePK:       v.NodePK,
			APR:          v.APY,
			VotePK:       v.VotePK,
			ActiveStake:  v.ActiveStake,
			Fee:          v.Fee,
			Score:        v.Score,
			SkippedSlots: v.SkippedSlots,
			DataCenter:   v.DataCenter,
		}
	}
	g.JSON(http.StatusOK, PoolDetails{
		Pool: Pool{
			Address:          pool.Address,
			Name:             pool.Address,
			ActiveStake:      pool.ActiveStake,
			TokensSupply:     pool.TokensSupply,
			APR:              pool.APR,
			AVGSkippedSlots:  pool.AVGSkippedSlots,
			AVGScore:         pool.AVGScore,
			Delinquent:       pool.Delinquent,
			UnstakeLiquidity: pool.UnstakeLiquidity,
			DepossitFee:      pool.DepossitFee,
			WithdrawalFee:    pool.WithdrawalFee,
			RewardsFee:       pool.RewardsFee,
		},
		Validators: validators,
	})
}
