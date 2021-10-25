package httpserv

import (
	"fmt"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	API struct {
		cfg config.Env
		svc services.Service
		log *zap.Logger
	}
)

func MewAPI(cfg config.Env, svc services.Service, log *zap.Logger) (api *API, err error) {
	return &API{
		cfg: cfg,
		svc: svc,
		log: log,
	}, nil
}

func (api API) Run() error {
	router := gin.Default()
	api.log.Info("Starting API server", zap.Uint64("port", api.cfg.HttpPort))
	return router.Run(fmt.Sprintf(":%d", api.cfg.HttpPort))
}
