package handler

import (
	"strings"

	"frame/internal/middleware"
	"frame/internal/module/admin/service"
	"frame/internal/pkg/cache"
	"frame/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{svc: service.NewAuthService()}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 10001, "参数错误: "+err.Error())
		return
	}

	res, err := h.svc.Login(&req)
	if err != nil {
		response.Fail(c, 20003, err.Error())
		return
	}

	response.OK(c, res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 {
		h.svc.Logout(parts[1])
	}

	userID := middleware.GetUserID(c)
	cache.ClearUserPermissions(userID)

	response.OK(c, nil)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userSvc := service.NewUserService()
	user, err := userSvc.GetProfile(userID)
	if err != nil {
		response.Fail(c, 20002, err.Error())
		return
	}
	response.OK(c, user)
}

func (h *AuthHandler) GetPermissions(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// Try cache first
	if perms, ok := cache.GetUserPermissions(userID); ok {
		response.OK(c, perms)
		return
	}

	menuSvc := service.NewMenuService()
	perms, err := menuSvc.GetUserPermissions(userID)
	if err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}

	cache.SetUserPermissions(userID, perms)
	response.OK(c, perms)
}
