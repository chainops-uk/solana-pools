package httpserv

import (
	"fmt"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/docs"
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	v1 "github.com/everstake/solana-pools/internal/delivery/httpserv/v1"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type (
	API struct {
		cfg config.Env
		svc services.Service
		log *zap.Logger
		v1  *v1.Handler
	}
)

func NewAPI(cfg config.Env, svc services.Service, log *zap.Logger) (api *API, err error) {
	return &API{
		cfg: cfg,
		svc: svc,
		log: log,
		v1:  v1.New(svc, log),
	}, nil
}

func (api API) Run() error {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/v1"
	v1g := router.Group("/v1")
	v1g.GET("/pools", tools.Must(api.v1.GetPools))
	v1g.GET("/pool/:name", tools.Must(api.v1.GetPool))
	v1g.GET("/pool-statistic", tools.Must(api.v1.GetPoolsStatistic))
	v1g.GET("/pools-statistic", tools.Must(api.v1.GetTotalPoolsStatistic))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	api.log.Info("Starting API server", zap.Uint64("port", api.cfg.HttpPort))
	return router.Run(fmt.Sprintf(":%d", api.cfg.HttpPort))
}
