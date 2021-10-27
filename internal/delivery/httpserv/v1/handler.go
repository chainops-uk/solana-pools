package v1

import (
	"github.com/everstake/solana-pools/internal/services"
	"go.uber.org/zap"
)

// @BasePath /v1

type Handler struct {
	svc services.Service
	log *zap.Logger
}

func New(svc services.Service, log *zap.Logger) *Handler {
	return &Handler{
		svc: svc,
		log: log,
	}
}
