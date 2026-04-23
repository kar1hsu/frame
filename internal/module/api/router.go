package api

import (
	"frame/internal/module/api/handler"
	"github.com/gin-gonic/gin"
)

type Module struct{}

func New() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return "api"
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	healthHandler := handler.NewHealthHandler()

	rg.GET("/health", healthHandler.Health)
}
