package service

import (
	"errors"
	"fmt"

	"github.com/karlhsu/frame/internal/app"
	"github.com/karlhsu/frame/internal/dao"
	"github.com/karlhsu/frame/internal/model"
	"github.com/karlhsu/frame/internal/pkg/cache"
)

type RoleService struct {
	roleDAO *dao.RoleDAO
}

func NewRoleService() *RoleService {
	return &RoleService{roleDAO: dao.NewRoleDAO()}
}

type CreateRoleRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

type UpdateRoleRequest struct {
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status"`
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
	if _, err := s.roleDAO.GetByCode(req.Code); err == nil {
		return errors.New("角色编码已存在")
	}

	role := &model.SysRole{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}
	return s.roleDAO.Create(role)
}

func (s *RoleService) GetByID(id uint) (*model.SysRole, error) {
	return s.roleDAO.GetByID(id)
}

func (s *RoleService) Update(id uint, req *UpdateRoleRequest) error {
	role, err := s.roleDAO.GetByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	role.Sort = req.Sort
	if req.Status != 0 {
		role.Status = req.Status
	}
	if req.Remark != "" {
		role.Remark = req.Remark
	}
	return s.roleDAO.Update(role)
}

func (s *RoleService) Delete(id uint) error {
	role, err := s.roleDAO.GetByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	// Remove all casbin policies for this role
	app.Enforcer.RemoveFilteredPolicy(0, role.Code)

	return s.roleDAO.Delete(id)
}

func (s *RoleService) List(page, pageSize int) ([]model.SysRole, int64, error) {
	return s.roleDAO.List(page, pageSize)
}

func (s *RoleService) ListAll() ([]model.SysRole, error) {
	return s.roleDAO.ListAll()
}

// permissionAPIs maps menu permission identifiers to API routes.
// When menus are assigned to a role, Casbin policies are auto-generated from this map.
func (s *RoleService) SetMenus(roleID uint, menuIDs []uint) error {
	if err := s.roleDAO.SetMenus(roleID, menuIDs); err != nil {
		return err
	}

	role, err := s.roleDAO.GetByID(roleID)
	if err != nil {
		return err
	}

	// Auto-sync Casbin policies from menu-associated APIs
	app.Enforcer.RemoveFilteredPolicy(0, role.Code)

	if len(menuIDs) > 0 {
		menuDAO := dao.NewMenuDAO()
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
	role, err := s.roleDAO.GetByID(roleID)
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
	role, err := s.roleDAO.GetByID(roleID)
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
