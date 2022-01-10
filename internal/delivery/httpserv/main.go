package httpserv

import (
	"fmt"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/docs"
	"github.com/everstake/solana-pools/internal/delivery/httpserv/tools"
	v1 "github.com/everstake/solana-pools/internal/delivery/httpserv/v1"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
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

func (api *API) Run() error {
	gin.SetMode(api.cfg.GinMode)

	router := gin.New()
	router.Use(ginzap.Ginzap(
		api.log, time.RFC3339, true),
		gin.Recovery(),
		cors.Default(),
	)

	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Title = "Stake-Solana API"
	docs.SwaggerInfo.Description = "This is API for Rest & WebSocket requests"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = api.cfg.HttpSwaggerAddress
	docs.SwaggerInfo.Schemes = []string{"http", "https", "ws"}
	router.GET("/", func(ctx *gin.Context) {
		links := []string{
			fmt.Sprintf(`<a href="http://%s">%s</a> - swagger`,
				api.cfg.HttpSwaggerAddress+"/swagger/index.html", api.cfg.HttpSwaggerAddress+"/swagger/v1/index.html"),
		}

		content := strings.Join(links, "\n")

		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.String(http.StatusCreated, `<html>%s</html>`, content)
	})

	v1g := router.Group("/v1")
	v1g.GET("/epoch", tools.Must(api.v1.GetEpoch))
	v1g.GET("/pools", tools.Must(api.v1.GetPools))
	v1g.GET("/coins", tools.Must(api.v1.GetCoins))
	v1g.GET("/pool-coins", tools.Must(api.v1.GetPoolsCoins))
	v1g.GET("/governance", tools.Must(api.v1.GetGovernance))
	v1g.GET("/pool-validators/:name", tools.Must(api.v1.GetPoolValidators))
	v1g.GET("/pool/:name", tools.WSMust(api.v1.GetPool, time.Second*30))
	v1g.GET("/pool-statistic", tools.Must(api.v1.GetPoolsStatistic))
	v1g.GET("/pools-statistic", tools.WSMust(api.v1.GetTotalPoolsStatistic, time.Second*30))
	v1g.GET("/liquidity-pool", tools.Must(api.v1.GetLiquidityPools))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	api.log.Info("Starting API server", zap.Uint64("port", api.cfg.HttpPort))
	return router.Run(fmt.Sprintf(":%d", api.cfg.HttpPort))
}
