package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/karlhsu/frame/internal/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	response.OK(c, gin.H{
		"status": "ok",
	})
}
