package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kar1hsu/frame/internal/middleware"
	"github.com/kar1hsu/frame/internal/module/admin/service"
	"github.com/kar1hsu/frame/internal/pkg/cache"
	"github.com/kar1hsu/frame/internal/pkg/errcode"
	"github.com/kar1hsu/frame/internal/pkg/response"
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
		response.Fail(c, errcode.ErrParam, "参数错误: "+err.Error())
		return
	}

	res, err := h.svc.Login(&req, c.ClientIP())
	if err != nil {
		response.Fail(c, errcode.ErrPasswordWrong, err.Error())
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
		response.Fail(c, errcode.ErrUserNotFound, err.Error())
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
		response.Fail(c, errcode.ErrServer, err.Error())
		return
	}

	cache.SetUserPermissions(userID, perms)
	response.OK(c, perms)
}
