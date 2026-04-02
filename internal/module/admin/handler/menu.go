package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/karlhsu/frame/internal/middleware"
	"github.com/karlhsu/frame/internal/module/admin/service"
	"github.com/karlhsu/frame/internal/pkg/response"
)

type MenuHandler struct {
	svc *service.MenuService
}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{svc: service.NewMenuService()}
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req service.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 10001, "参数错误: "+err.Error())
		return
	}
	if req.Status == 0 {
		req.Status = 1
	}
	if req.Visible == 0 {
		req.Visible = 1
	}

	if err := h.svc.Create(&req); err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *MenuHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	menu, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, 50001, "菜单不存在")
		return
	}
	response.OK(c, menu)
}

func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	var req service.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 10001, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.Update(uint(id), &req); err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *MenuHandler) GetTree(c *gin.Context) {
	tree, err := h.svc.GetTree()
	if err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, tree)
}

func (h *MenuHandler) GetUserMenuTree(c *gin.Context) {
	userID := middleware.GetUserID(c)
	tree, err := h.svc.GetUserMenuTree(userID)
	if err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, tree)
}
