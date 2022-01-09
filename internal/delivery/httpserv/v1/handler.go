package v1

import (
	"github.com/everstake/solana-pools/internal/services"
	"go.uber.org/zap"
)

//go:generate swag init -g ./handler.go -d ./,../tools -o ../../../../docs

// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1
// @query.collection.format multi

// @x-extension-openapi {"example": "value on a json format"}

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
