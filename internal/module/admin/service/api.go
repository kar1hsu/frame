package service

import (
	"github.com/kar1hsu/frame/internal/model"
	"github.com/kar1hsu/frame/internal/repository"
)

type APIService struct {
	apiRepo *repository.ApiRepo
}

func NewAPIService() *APIService {
	return &APIService{apiRepo: repository.NewApiRepo()}
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
	return s.apiRepo.Create(api)
}

func (s *APIService) GetByID(id uint) (*model.SysAPI, error) {
	return s.apiRepo.GetByID(id)
}

func (s *APIService) Update(id uint, req *UpdateAPIRequest) error {
	api, err := s.apiRepo.GetByID(id)
	if err != nil {
		return notFoundOr(err, "API 不存在")
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
	return s.apiRepo.Update(api)
}

func (s *APIService) Delete(id uint) error {
	return s.apiRepo.Delete(id)
}

func (s *APIService) List(page, pageSize int) ([]model.SysAPI, int64, error) {
	return s.apiRepo.List(page, pageSize)
}

func (s *APIService) ListAll() ([]model.SysAPI, error) {
	return s.apiRepo.ListAll()
}
