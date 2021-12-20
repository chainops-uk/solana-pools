package v1

import (
	"errors"
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// GetPoolValidators godoc
// @Summary REST
// @Schemes
// @Description get pool validators
// @Param name path string true "Pool name" default(marinade)
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Accept json
// @Produce json
// @Success 200 {object} tools.ResponseArrayData{data=[]validator} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pool-validators/{name} [get]
func (h *Handler) GetPoolValidators(ctx *gin.Context) (interface{}, error) {
	name := ctx.Param("name")
	q := struct {
		Name   string `form:"name"`
		Offset uint64 `form:"offset,default=0"`
		Limit  uint64 `form:"limit,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	resp, count, err := h.svc.GetPoolValidators(name, q.Limit, q.Offset)
	if err != nil {
		h.log.Error("API GetPoolData", zap.Error(err))
		if errors.Is(err, postgres.ErrorRecordNotFounded) {
			return nil, tools.NewStatus(http.StatusBadRequest, fmt.Errorf("%s pool not found", name))
		}
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	arr := make([]*validator, len(resp))
	for i, v := range resp {
		arr[i] = (&validator{}).Set(v)
	}

	return tools.ResponseArrayData{
		Data: arr,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: count,
		},
	}, nil
}

type validator struct {
	Name             string  `json:"name"`
	Image            string  `json:"image"`
	NodePK           string  `json:"node_pk"`
	APY              float64 `json:"apy"`
	VotePK           string  `json:"vote_pk"`
	PoolActiveStake  float64 `json:"pool_active_stake"`
	TotalActiveStake float64 `json:"total_active_stake"`
	Fee              float64 `json:"fee"`
	Score            int64   `json:"score"`
	SkippedSlots     float64 `json:"skipped_slots"`
	DataCenter       string  `json:"data_center"`
}

func (v *validator) Set(validator *smodels.Validator) *validator {
	v.NodePK = validator.NodePK
	v.Name = validator.Name
	v.Image = validator.Image
	v.APY, _ = validator.APY.Float64()
	v.VotePK = validator.VotePK
	v.PoolActiveStake, _ = validator.PoolActiveStake.Float64()
	v.TotalActiveStake, _ = validator.TotalActiveStake.Float64()
	v.Fee, _ = validator.Fee.Float64()
	v.Score = validator.Score
	v.SkippedSlots, _ = validator.SkippedSlots.Float64()
	v.DataCenter = validator.DataCenter

	return v
}