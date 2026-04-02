package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/karlhsu/frame/internal/module/admin/service"
	"github.com/karlhsu/frame/internal/pkg/response"
	"github.com/karlhsu/frame/internal/pkg/utils"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{svc: service.NewUserService()}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req service.CreateUserRequest
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

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	user, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, 20002, "用户不存在")
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}

	var req service.UpdateUserRequest
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

func (h *UserHandler) Delete(c *gin.Context) {
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

func (h *UserHandler) List(c *gin.Context) {
	p := utils.GetPagination(c)
	users, total, err := h.svc.List(p.Page, p.PageSize)
	if err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, gin.H{
		"list":  users,
		"total": total,
		"page":  p.Page,
		"size":  p.PageSize,
	})
}
