package v1

import (
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetCoins godoc
// @Summary RestAPI
// @Schemes
// @Description get coins
// @Accept json
// @Produce json
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Param name query string false "stake-pool name"
// @Success 200 {object} tools.ResponseData{data=[]poolMainPage} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /coins [get]
func (h *Handler) GetCoins(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string `form:"name"`
		Limit  uint64 `form:"limit,default=1"`
		Offset uint64 `form:"offset,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	scoins, count, err := h.svc.GetCoins(q.Name, q.Limit, q.Offset)
	if err != nil {
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	coins := make([]*coin, len(scoins))
	for i, c := range scoins {
		coins[i] = (&coin{}).Set(c)
	}

	return tools.ResponseArrayData{
		Data: scoins,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: count,
		},
	}, nil
}

// GetPoolsCoins godoc
// @Summary RestAPI
// @Schemes
// @Description get pools coins
// @Accept json
// @Produce json
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Param name query string false "stake-pool name"
// @Success 200 {object} tools.ResponseData{data=[]coin} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /pool-coins [get]
func (h *Handler) GetPoolsCoins(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string `form:"name"`
		Offset uint64 `form:"offset,default=0"`
		Limit  uint64 `form:"limit,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	scoins, count, err := h.svc.GetPoolCoins(q.Name, q.Limit, q.Offset)
	if err != nil {
		return nil, tools.NewStatus(http.StatusInternalServerError, err)
	}

	coins := make([]*coin, len(scoins))
	for i, c := range scoins {
		coins[i] = (&coin{}).Set(c)
	}

	return tools.ResponseArrayData{
		Data: coins,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: count,
		},
	}, nil
}

type coin struct {
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	USD        float64 `json:"usd"`
	ThumbImage string  `json:"thumb_image"`
	SmallImage string  `json:"small_image"`
	LargeImage string  `json:"large_image"`
}

func (c *coin) Set(coin *smodels.Coin) *coin {
	c.USD = coin.USD
	c.ThumbImage = coin.ThumbImage
	c.SmallImage = coin.SmallImage
	c.LargeImage = coin.LargeImage
	c.Name = coin.Name
	c.Address = coin.Address
	return c
}
