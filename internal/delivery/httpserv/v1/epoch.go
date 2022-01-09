package v1

import (
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	"github.com/gin-gonic/gin"
)

// GetEpoch godoc
// @Summary RestAPI
// @Schemes
// @Description The current epoch value is returned.
// @Tags epoch
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
