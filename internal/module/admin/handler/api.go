package handler

import (
	"strconv"

	"frame/internal/module/admin/service"
	"frame/internal/pkg/response"
	"frame/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	svc *service.APIService
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{svc: service.NewAPIService()}
}

func (h *APIHandler) Create(c *gin.Context) {
	var req service.CreateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 10001, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Create(&req); err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *APIHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}
	api, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, 10002, "API 不存在")
		return
	}
	response.OK(c, api)
}

func (h *APIHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 10001, "参数错误")
		return
	}
	var req service.UpdateAPIRequest
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

func (h *APIHandler) Delete(c *gin.Context) {
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

func (h *APIHandler) List(c *gin.Context) {
	p := utils.GetPagination(c)
	apis, total, err := h.svc.List(p.Page, p.PageSize)
	if err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, gin.H{
		"list":  apis,
		"total": total,
		"page":  p.Page,
		"size":  p.PageSize,
	})
}

func (h *APIHandler) ListAll(c *gin.Context) {
	apis, err := h.svc.ListAll()
	if err != nil {
		response.Fail(c, 10000, err.Error())
		return
	}
	response.OK(c, apis)
}
