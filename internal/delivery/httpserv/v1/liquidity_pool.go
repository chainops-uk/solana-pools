package v1

import (
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// GetLiquidityPools godoc
// @Summary RestAPI
// @Schemes
// @Description This Liquidity Pools list with search by name.
// @Tags pool
// @Accept json
// @Produce json
// @Param name query string false "The name of the pool without strict observance of the case."
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Success 200 {object} tools.ResponseArrayData{data=[]liquidityPool} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /liquidity-pools [get]
func (h *Handler) GetLiquidityPools(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string `form:"name"`
		Offset uint64 `form:"offset,default=0"`
		Limit  uint64 `form:"limit,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	pools, amount, err := h.svc.GetLiquidityPools(q.Name, q.Limit, q.Offset)
	if err != nil {
		h.log.Error("API GetPoolData", zap.Error(err))
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	aPools := make([]*liquidityPool, len(pools))
	for i, v := range pools {
		aPools[i] = (&liquidityPool{}).Set(v)
	}

	return tools.ResponseArrayData{
		Data: aPools,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: amount,
		}}, nil
}
