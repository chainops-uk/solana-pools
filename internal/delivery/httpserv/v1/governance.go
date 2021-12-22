package v1

import (
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
)

// GetGovernance godoc
// @Summary RestAPI
// @Schemes
// @Description get governance
// @Accept json
// @Produce json
// @Param offset query number true "offset for aggregation" default(0)
// @Param limit query number true "limit for aggregation" default(10)
// @Param name query string false "stake-pool name"
// @Success 200 {object} tools.ResponseData{data=[]governance} "Ok"
// @Failure 400,404 {object} tools.ResponseError "bad request"
// @Failure 500 {object} tools.ResponseError "internal server error"
// @Failure default {object} tools.ResponseError "default response"
// @Router /governance [get]
func (h *Handler) GetGovernance(ctx *gin.Context) (interface{}, error) {
	q := struct {
		Name   string `form:"name"`
		Offset uint64 `form:"offset,default=0"`
		Limit  uint64 `form:"limit,default=10"`
	}{}
	if err := ctx.ShouldBind(&q); err != nil {
		return nil, tools.NewStatus(http.StatusBadRequest, err)
	}

	gc, count, err := h.svc.GetGovernance(q.Name, q.Limit, q.Offset)
	if err != nil {
		return nil, err
	}

	g := make([]*governance, len(gc))
	for i, s := range gc {
		g[i] = (&governance{}).Set(s)
	}

	return tools.ResponseArrayData{
		Data: g,
		MetaData: &tools.MetaData{
			Offset:      q.Offset,
			Limit:       q.Limit,
			TotalAmount: count,
		},
	}, err
}

type governance struct {
	Name                  string          `json:"name"`
	Image                 string          `json:"image"`
	Blockchain            string          `json:"blockchain"`
	ContractAddress       string          `json:"contract_address"`
	VoteURL               string          `json:"vote_url"`
	About                 string          `json:"about"`
	Vote                  string          `json:"vote"`
	Trade                 string          `json:"trade"`
	Exchange              string          `json:"exchange"`
	MaximumTokenSupply    float64         `json:"maximum_token_supply"`
	CirculatingSupply     float64         `json:"circulating_supply"`
	USD                   float64         `json:"usd"`
	DAOTreasury           decimal.Decimal `json:"dao_treasury"`
	Investors             decimal.Decimal `json:"investors"`
	InitialLidoDevelopers decimal.Decimal `json:"initial_lido_developers"`
	Foundation            decimal.Decimal `json:"foundation"`
	Validators            decimal.Decimal `json:"validators"`
}

func (g *governance) Set(governance *smodels.Governance) *governance {
	g.Name = governance.Name
	g.Image = governance.Image
	g.Blockchain = governance.Blockchain
	g.ContractAddress = governance.ContractAddress
	g.VoteURL = governance.VoteURL
	g.About = governance.About
	g.Vote = governance.Vote
	g.Trade = governance.Trade
	g.Exchange = governance.Exchange
	g.MaximumTokenSupply = governance.MaximumTokenSupply
	g.CirculatingSupply = governance.CirculatingSupply
	g.USD = governance.USD
	g.DAOTreasury = governance.DAOTreasury
	g.Investors = governance.Investors
	g.InitialLidoDevelopers = governance.InitialLidoDevelopers
	g.Foundation = governance.Foundation
	g.Validators = governance.Validators
	return g
}
