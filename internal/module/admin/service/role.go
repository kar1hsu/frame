package service

import (
	"errors"
	"fmt"

	"github.com/kar1hsu/frame/internal/app"
	"github.com/kar1hsu/frame/internal/model"
	"github.com/kar1hsu/frame/internal/pkg/cache"
	"github.com/kar1hsu/frame/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepo
}

func NewRoleService() *RoleService {
	return &RoleService{roleRepo: repository.NewRoleRepo()}
}

type CreateRoleRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Sort   int    `json:"sort"`
	Status *int8  `json:"status"`
	Remark string `json:"remark"`
}

type UpdateRoleRequest struct {
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
	Status *int8  `json:"status"`
	Remark string `json:"remark"`
}

type SetRoleMenusRequest struct {
	MenuIDs []uint `json:"menu_ids" binding:"required"`
}

type SetRoleAPIsRequest struct {
	APIs []RoleAPIItem `json:"apis" binding:"required"`
}

type RoleAPIItem struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

func (s *RoleService) Create(req *CreateRoleRequest) error {
	if _, err := s.roleRepo.GetByCode(req.Code); err == nil {
		return errors.New("角色编码已存在")
	}

	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	role := &model.SysRole{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: status,
		Remark: req.Remark,
	}
	return s.roleRepo.Create(role)
}

func (s *RoleService) GetByID(id uint) (*model.SysRole, error) {
	return s.roleRepo.GetByID(id)
}

func (s *RoleService) Update(id uint, req *UpdateRoleRequest) error {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	role.Sort = req.Sort
	if req.Status != nil {
		role.Status = *req.Status
	}
	if req.Remark != "" {
		role.Remark = req.Remark
	}
	return s.roleRepo.Update(role)
}

func (s *RoleService) Delete(id uint) error {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	// Remove all casbin policies for this role
	app.Enforcer.RemoveFilteredPolicy(0, role.Code)

	return s.roleRepo.Delete(id)
}

func (s *RoleService) List(page, pageSize int) ([]model.SysRole, int64, error) {
	return s.roleRepo.List(page, pageSize)
}

func (s *RoleService) ListAll() ([]model.SysRole, error) {
	return s.roleRepo.ListAll()
}

// permissionAPIs maps menu permission identifiers to API routes.
// When menus are assigned to a role, Casbin policies are auto-generated from this map.
func (s *RoleService) SetMenus(roleID uint, menuIDs []uint) error {
	if err := s.roleRepo.SetMenus(roleID, menuIDs); err != nil {
		return err
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	// Auto-sync Casbin policies from menu-associated APIs
	app.Enforcer.RemoveFilteredPolicy(0, role.Code)

	if len(menuIDs) > 0 {
		menuDAO := repository.NewMenuRepo()
		menus, err := menuDAO.GetByIDs(menuIDs)
		if err != nil {
			return err
		}

		for _, m := range menus {
			for _, api := range m.APIs {
				app.Enforcer.AddPolicy(role.Code, api.Path, api.Method)
			}
		}
	}

	cache.ClearAllPermissionCache()

	return app.Enforcer.SavePolicy()
}

func (s *RoleService) SetAPIs(roleID uint, apis []RoleAPIItem) error {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// Clear old policies
	app.Enforcer.RemoveFilteredPolicy(0, role.Code)

	// Add new policies
	for _, api := range apis {
		if _, err := app.Enforcer.AddPolicy(role.Code, api.Path, api.Method); err != nil {
			return fmt.Errorf("添加权限策略失败: %w", err)
		}
	}

	return app.Enforcer.SavePolicy()
}

func (s *RoleService) GetAPIs(roleID uint) ([]RoleAPIItem, error) {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	policies := app.Enforcer.GetFilteredPolicy(0, role.Code)
	items := make([]RoleAPIItem, 0, len(policies))
	for _, p := range policies {
		if len(p) >= 3 {
			items = append(items, RoleAPIItem{Path: p[1], Method: p[2]})
		}
	}
	return items, nil
}
