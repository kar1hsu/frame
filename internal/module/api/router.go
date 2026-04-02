package api

import (
	"github.com/gin-gonic/gin"
	"github.com/karlhsu/frame/internal/module/api/handler"
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
