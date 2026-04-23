package handler

import (
	"strconv"

	"frame/internal/module/admin/service"
	"frame/internal/pkg/response"
	"frame/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{svc: service.NewRoleService()}
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 10001, "参数错误: "+err.Error())
		return
	}
	if req.Status == 0 {
		req.Status = 1
	}

	if err := h.svc.Create(&req); err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	role, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, 40002, "角色不存在")
		return
	}
	response.OK(c, role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	var req service.UpdateRoleRequest
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

func (h *RoleHandler) Delete(c *gin.Context) {
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

func (h *RoleHandler) List(c *gin.Context) {
	p := utils.GetPagination(c)
	roles, total, err := h.svc.List(p.Page, p.PageSize)
	if err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, gin.H{
		"list":  roles,
		"total": total,
		"page":  p.Page,
		"size":  p.PageSize,
	})
}

func (h *RoleHandler) ListAll(c *gin.Context) {
	roles, err := h.svc.ListAll()
	if err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, roles)
}

func (h *RoleHandler) SetMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	var req service.SetRoleMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 10001, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.SetMenus(uint(id), req.MenuIDs); err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RoleHandler) SetAPIs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	var req service.SetRoleAPIsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 10001, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.SetAPIs(uint(id), req.APIs); err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *RoleHandler) GetAPIs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	apis, err := h.svc.GetAPIs(uint(id))
	if err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, apis)
}
