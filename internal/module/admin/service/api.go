package service

import (
	"errors"

	"frame/internal/dao"
	"frame/internal/model"
)

type APIService struct {
	apiDAO *dao.APIDAO
}

func NewAPIService() *APIService {
	return &APIService{apiDAO: dao.NewAPIDAO()}
}

type CreateAPIRequest struct {
	Path        string `json:"path" binding:"required"`
	Method      string `json:"method" binding:"required"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

type UpdateAPIRequest struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

func (s *APIService) Create(req *CreateAPIRequest) error {
	api := &model.SysAPI{
		Path:        req.Path,
		Method:      req.Method,
		Group:       req.Group,
		Description: req.Description,
	}
	return s.apiDAO.Create(api)
}

func (s *APIService) GetByID(id uint) (*model.SysAPI, error) {
	return s.apiDAO.GetByID(id)
}

func (s *APIService) Update(id uint, req *UpdateAPIRequest) error {
	api, err := s.apiDAO.GetByID(id)
	if err != nil {
		return errors.New("API 不存在")
	}
	if req.Path != "" {
		api.Path = req.Path
	}
	if req.Method != "" {
		api.Method = req.Method
	}
	if req.Group != "" {
		api.Group = req.Group
	}
	if req.Description != "" {
		api.Description = req.Description
	}
	return s.apiDAO.Update(api)
}

func (s *APIService) Delete(id uint) error {
	return s.apiDAO.Delete(id)
}

func (s *APIService) List(page, pageSize int) ([]model.SysAPI, int64, error) {
	return s.apiDAO.List(page, pageSize)
}

func (s *APIService) ListAll() ([]model.SysAPI, error) {
	return s.apiDAO.ListAll()
}
